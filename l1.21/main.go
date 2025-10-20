package main

import "fmt"

// Пусть у нас есть банк где по своему определены методы снятия/вывода и пополнения средств на балансе
type WildberriesBank struct {
}

// Таких банков может быть много и у каждого по своему определены пополнения и вывода средств
func (wb *WildberriesBank) Charge(req PaymentRequest) error {
	fmt.Println("Charged", req.Amount, ", Wildberries Bank with id", req.Id)
	return nil
}
func (wb *WildberriesBank) TopUp(amount int64) error {
	fmt.Println("Topped up", amount, "to Wildberries Bank")
	return nil
}

// Условный конструктор для платежа
type PaymentRequest struct {
	Amount int64
	Id     string
}

// Клиентский контракт
// Пусть это магазин, и за каждый заказ он может либо провести оплату (то есть списать деньги с точки зрения банка)
// либо вернуть деньги (то есть пополнить счет с точки зрения банка)
type PaymentGateway interface {
	Pay(req PaymentRequest) error
	Refund(amount int64) error
}

// Теперь для каждого банка нужен адаптер, то есть структура в которой будет целевой банк
// И одноименные методы из интерфейса, но функции работают исключительно с методами целевого банка
type WildberriesAdapter struct {
	bank *WildberriesBank
}

func (wa *WildberriesAdapter) Pay(req PaymentRequest) error {
	return wa.bank.Charge(req)
}
func (wa *WildberriesAdapter) Refund(amount int64) error {
	return wa.bank.TopUp(amount)
}

// Можно еще добавить много таких банков но для каждого нужно написать адаптер + переопределленые методы, но называться они должны как в интерфейсе

// Платеж, функция клиента (Магазина например)
func ProcessPayment(gateway PaymentGateway) {
	req := PaymentRequest{Amount: 1000, Id: "12345"}
	gateway.Pay(req)
	gateway.Refund(500)
}
func main() {
	bank := &WildberriesBank{}
	adapter := &WildberriesAdapter{bank: bank}
	ProcessPayment(adapter) //Клиент не знает о WildberriesBank
}
