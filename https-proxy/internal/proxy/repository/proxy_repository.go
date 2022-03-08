package repository

import (
	"context"
	"errors"

	"http-proxy/internal/proxy/models"
	"http-proxy/pkg/utils"
)

type ProxyRepository struct {
	conn *utils.PostgresConn
}

const (
	insertRequestQuery  = `INSERT INTO requests(method, path, get_params, headers, cookies, post_params) VALUES($1, $2, $3, $4, $5, $6) RETURNING id;`
	insertResponseQuery = `INSERT INTO responses(request_id, code, message, headers, body) VALUES($1, $2, $3, $4, $5);`
)

func NewProxyRepository(conn *utils.PostgresConn) *ProxyRepository {
	return &ProxyRepository{
		conn: conn,
	}
}
func (p *ProxyRepository) InsertRequest(req *models.Request) (int, error) {
	id := -1
	err := p.conn.Conn.QueryRow(context.Background(), insertRequestQuery, req.Method, req.Path, req.GetParams, req.Headers, req.Cookies, req.PostParams).Scan(&id)
	if err != nil {
		return id, err
	}
	return id, nil

}

func (p *ProxyRepository) InsertResponse(reqID int, resp *models.Response) error {
	res, err := p.conn.Conn.Exec(context.Background(), insertResponseQuery, reqID, resp.Code, resp.Message, resp.Headers, resp.Body)
	if err != nil {
		return err
	}
	if res.RowsAffected() != 1 {
		return errors.New("internal server error")
	}
	return nil
}
