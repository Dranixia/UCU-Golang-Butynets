package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Retriable interface {
	GetMaxAttempts() uint
	Retry(action func() bool)
}

type DefaultRetriable struct {
	maxRetries uint
}

type RetriableWithDelay struct {
	maxRetries uint
	delaySecs  time.Duration
}

func (dr DefaultRetriable) GetMaxAttempts() uint {
	return dr.maxRetries
}

func (rd RetriableWithDelay) GetMaxAttempts() uint {
	return rd.maxRetries
}

func (rd RetriableWithDelay) Retry(action func() bool) {
	attempts := rd.GetMaxAttempts()
	for i := 0; i < int(attempts); i++ {
		if action() == true {
			fmt.Printf("Success after attempt № %d\n", i+1)
			return
		}
		fmt.Printf("Failure after attempt № %d\n", i+1)
		time.Sleep(time.Second * rd.delaySecs)
	}
}

func (dr DefaultRetriable) Retry(action func() bool) {
	attempts := dr.GetMaxAttempts()
	for i := 0; i < int(attempts); i++ {
		if action() == true {
			fmt.Printf("Success after attempt № %d\n", i+1)
			return
		}
		fmt.Printf("Failure after attempt № %d\n", i+1)
	}
}

func NewDefaultRetriable(maxRetries ...uint) Retriable {
	var maxRetriesArg uint = 10
	if len(maxRetries) > 0 {
		maxRetriesArg = maxRetries[0]
	}
	return &DefaultRetriable{maxRetries: maxRetriesArg}
}

func NewRetriableWithDelay(delayInSecs uint, maxRetries ...uint) Retriable {
	var maxRetriesArg uint = 10
	if len(maxRetries) > 0 {
		maxRetriesArg = maxRetries[0]
	}
	return &RetriableWithDelay{maxRetries: maxRetriesArg, delaySecs: time.Duration(delayInSecs)}
}

// Реалізуйте інтерфейс у двох варіантах DefaultRetriable, RetriableWithDelay
// В обох структурах Retry має викликати передану функцію не більше ніж повертає GetMaxAttempts.
// Якщо action() поверне true - зупиняєте виконання циклу
// RetriableWithDelay структура має мати поле DelayInSec. У випадку якщо action() повернув false - програма має почекати вказану кількість секунд (time.Sleep)
// Якщо action() поверне false - вивести в консоль повідомлення в якому вказати номер неуспішної спроби
// приклад використання
//
//	ret := RetriableWithDelay{
//		DelayInSec: 1,
//
// fmt
//
//	ret.Retry(func() bool {
//		return false
//	})
func main() {
	var testFunc = func() bool {
		if val := rand.Intn(10); val == 0 {
			return true
		} else {
			return false
		}
	}

	rand.Seed(time.Now().UnixNano())

	retDefault := NewDefaultRetriable()
	fmt.Printf("DefaultRetriable, Attempts: %d\n", retDefault.GetMaxAttempts())
	retDefault.Retry(testFunc)

	retDefaultWithTenAttempts := NewDefaultRetriable(5)
	fmt.Printf("DefaultRetriable, Attempts: %d\n", retDefaultWithTenAttempts.GetMaxAttempts())
	retDefaultWithTenAttempts.Retry(testFunc)

	retDelay := NewRetriableWithDelay(1)
	fmt.Printf("RetriableWithDelay, Attempts: %d\n", retDelay.GetMaxAttempts())
	retDelay.Retry(testFunc)

	retDelayWithTenAttempts := NewRetriableWithDelay(1, 5)
	fmt.Printf("RetriableWithDelay, Attempts: %d\n", retDelayWithTenAttempts.GetMaxAttempts())
	retDelayWithTenAttempts.Retry(testFunc)
}
