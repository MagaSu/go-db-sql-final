package main

import (
	"database/sql"
	"fmt"
)

type ParcelStore struct {
	db *sql.DB
}

func NewParcelStore(db *sql.DB) ParcelStore {
	return ParcelStore{db: db}
}

func (s ParcelStore) Add(p Parcel) (int, error) {
	// реализуйте добавление строки в таблицу parcel, используйте данные из переменной p
	query := `INSERT INTO parcel (client, status, address, created_at)
	VALUES (:client, :status, :address, :created_at)`

	req, err := s.db.Exec(query,
		sql.Named("client", p.Client),
		sql.Named("status", p.Status),
		sql.Named("address", p.Address),
		sql.Named("created_at", p.CreatedAt),
	)
	if err != nil {
		return 0, err
	}
	lastId, err := req.LastInsertId()
	if err != nil {
		return 0, err
	}

	// верните идентификатор последней добавленной записи
	return int(lastId), nil
}

func (s ParcelStore) Get(number int) (Parcel, error) {
	// реализуйте чтение строки по заданному number
	// здесь из таблицы должна вернуться только одна строка
	var (
		client    int
		status    string
		address   string
		createdAt string
	)

	query := `SELECT client, status, address, created_at FROM parcel WHERE number = :number`
	parcel, err := s.db.Query(query, sql.Named("number", number))
	if err != nil {
		return Parcel{}, err
	}
	for parcel.Next() {

		err = parcel.Scan(&client, &status, &address, &createdAt)
		if err != nil {
			return Parcel{}, err
		}
	}
	if err := parcel.Err(); err != nil {
		return Parcel{}, err
	}

	// заполните объект Parcel данными из таблицы
	p := Parcel{
		Number:    number,
		Client:    client,
		Status:    status,
		Address:   address,
		CreatedAt: createdAt,
	}

	return p, nil
}

func (s ParcelStore) GetByClient(client int) ([]Parcel, error) {
	// реализуйте чтение строк из таблицы parcel по заданному client
	// здесь из таблицы может вернуться несколько строк
	var res []Parcel

	query := `SELECT number, status, address, created_at FROM parcel WHERE client = :client`
	rows, err := s.db.Query(query, sql.Named("client", client))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {

		var (
			number    int
			status    string
			address   string
			createdAt string
		)

		err = rows.Scan(&number, &status, &address, &createdAt)

		p := Parcel{
			Number:    number,
			Client:    client,
			Status:    status,
			Address:   address,
			CreatedAt: createdAt,
		}

		res = append(res, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return res, nil
}

func (s ParcelStore) SetStatus(number int, status string) error {
	// реализуйте обновление статуса в таблице parcel
	query := `UPDATE parcel SET status = :status WHERE number = :number`
	_, err := s.db.Exec(query,
		sql.Named("status", status),
		sql.Named("number", number),
	)
	if err != nil {
		return err
	}

	return nil
}

func (s ParcelStore) SetAddress(number int, address string) error {
	// реализуйте обновление адреса в таблице parcel
	// менять адрес можно только если значение статуса registered
	var status string
	query := `SELECT status FROM parcel WHERE number = :number`

	row := s.db.QueryRow(query, sql.Named("number", number))
	err := row.Scan(&status)
	if err == sql.ErrNoRows {
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("Failed to change the address: parcel not registered.")
	}

	query = `UPDATE parcel SET address = :address WHERE number = :number`
	_, err = s.db.Exec(query,
		sql.Named("address", address),
		sql.Named("number", number),
	)

	return nil
}

func (s ParcelStore) Delete(number int) error {
	// реализуйте удаление строки из таблицы parcel
	// удалять строку можно только если значение статуса registered
	var status string
	query := `SELECT status FROM parcel WHERE number = :number`

	row := s.db.QueryRow(query, sql.Named("number", number))
	err := row.Scan(&status)
	if err == sql.ErrNoRows {
		return err
	}

	if status != ParcelStatusRegistered {
		return fmt.Errorf("Failed to delete the parcel: parcel not registered.")
	}

	query = `DELETE FROM parcel WHERE number = :number`
	_, err = s.db.Exec(query, sql.Named("number", number))

	return nil
}
