package postgres

import (
	_ "github.com/lib/pq"
	"bitbucket.org/hornetdefiant/core-engine/factory"
)

func (p *Postgres) GetTags() ([]factory.Tag, error) {
	query := `
        SELECT value, counter
        FROM tags
        ORDER BY counter DESC
        LIMIT 5
    `

	// Execute the query and fetch the results
	rows, err := p.dbConn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]factory.Tag, 0)

	for rows.Next() {
		tag := factory.Tag{}
		if err := rows.Scan(&tag.Value, &tag.Counter); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}
