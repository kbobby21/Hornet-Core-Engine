package postgres

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (p *Postgres) AddWebSiteMeta(data factory.DarkWebSite) error {
	var err error

	// Eg:- insert into tags(value,counter) values('drugs',1),('child_porn',1),('arms',1) ON CONFLICT(value) DO UPDATE set counter=tags.counter+1 ;
	tagsQuery := "INSERT INTO tags(value) VALUES "
	tagsArgs := make([]interface{}, 0)
	for i, tag := range data.Tags {
		tagsQuery += fmt.Sprintf("($%d)", i+1)
		if i < len(data.Tags)-1 {
			tagsQuery += ","
		}
		tagsArgs = append(tagsArgs, tag)
	}

	tagsQuery += " ON CONFLICT(value) DO UPDATE set counter=tags.counter+1"

	if len(data.Tags) != 0 {
		_, err := p.dbConn.Exec(tagsQuery, tagsArgs...)
		if err != nil {
			return fmt.Errorf("error in inserting to tags table: %s", err)
		}
	}

	walletQuery := "INSERT INTO flagged_wallets (address, crypto_type) VALUES "

	walletQueryArgs := make([]interface{}, 0)
	for i, wallet := range data.Wallets {
		walletQuery += fmt.Sprintf("($%d,$%d)", 2*i+1, 2*i+2)
		if i < len(data.Wallets)-1 {
			walletQuery += ","
		}
		walletQueryArgs = append(walletQueryArgs, wallet.Address, wallet.CryptoType)
	}
	walletQuery += " ON CONFLICT DO NOTHING"

	if len(data.Wallets) != 0 {
		_, err = p.dbConn.Exec(
			walletQuery,
			walletQueryArgs...,
		)
		if err != nil {
			return fmt.Errorf("error in inserting to flagged_wallets table: %s", err)
		}
	}

	// join the list of strings of the body into one string by
	// joining on comma separator
	bodyTxt := strings.Join(data.Body, ",")

	dwQuery := "INSERT INTO dark_web_sites (website_name, onion_url, wallet_address, tag, body) VALUES "
	dwArgs := make([]interface{}, 0)
	c := 0
	lT := len(data.Tags)    // number of tags for this website
	lW := len(data.Wallets) // number of wallet addresses found on this website

	if lT >= lW {
		c = lT
	} else {
		c = lW
	}

	tempWallet := ""
	tempTag := ""
	for i := 0; i < c; i++ {
		// if number of wallets is exhausted use prev valid value
		if i < lW {
			tempWallet = data.Wallets[i].Address
		}

		if i < lT {
			tempTag = data.Tags[i]
		}

		dwQuery += fmt.Sprintf(
			"($%d,$%d,$%d,$%d,$%d)",
			5*i+1,
			5*i+2,
			5*i+3,
			5*i+4,
			5*i+5,
		)
		if i < c-1 {
			dwQuery += ","
		}
		dwArgs = append(
			dwArgs,
			data.WebsiteName,
			data.OnionUrl,
			tempWallet,
			tempTag,
			bodyTxt,
		)
	}

	_, err = p.dbConn.Exec(
		dwQuery,
		dwArgs...,
	)
	if err != nil {
		return fmt.Errorf("error in inserting to dark_web_sites table: %s", err)
	}

	return err
}

func (p *Postgres) GetWebsiteMeta(pageNum int) ([]factory.DarkwebResponse, error) {

	pageSize := viper.GetInt("page_size")
	offset := (pageNum - 1) * pageSize

	query := `SELECT
			website_name,
			onion_url,
			wallet_address,
			tag,
			body,
			discovered_at
		 FROM dark_web_sites 
		 ORDER BY discovered_at DESC LIMIT $1 OFFSET $2`

	rows, err := p.dbConn.Query(query, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("error in getting darkweb details: %s", err)
	}
	defer rows.Close()
	dw := make([]factory.DarkwebResponse, 0)
	for rows.Next() {
		var wn, ou, wa, t, b, da string
		if err := rows.Scan(
			&wn,
			&ou,
			&wa,
			&t,
			&b,
			&da,
		); err != nil {
			return nil, fmt.Errorf("error in scanning rows: %s", err)
		}
		var d factory.DarkwebResponse
		d.WebsiteName = wn
		d.OnionUrl = ou
		d.WalletAddress = wa
		d.Tag = t
		d.Body = b
		d.DiscoveredAt = da
		dw = append(dw, d)
	}
	return dw, nil
}
