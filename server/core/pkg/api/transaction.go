package api

import "github.com/sirupsen/logrus"

type TransactionStep interface {
	Execute() error
	Rollback() error
}

type Function func() error

type SimpleTransactionStep struct {
	ExecuteFunc  Function
	RollbackFunc Function
}

func (sts SimpleTransactionStep) Execute() error {
	if sts.ExecuteFunc != nil {
		return sts.ExecuteFunc()
	}

	return nil
}

func (sts SimpleTransactionStep) Rollback() error {
	if sts.RollbackFunc != nil {
		return sts.ExecuteFunc()
	}

	return nil
}

type Transaction struct {
	steps []TransactionStep
}

func (tx *Transaction) AddStep(step TransactionStep) {
	tx.steps = append(tx.steps, step)
}

func (tx *Transaction) AddFunction(callback Function) {
	tx.steps = append(tx.steps, SimpleTransactionStep{
		ExecuteFunc: callback,
	})
}

func (tx *Transaction) AddRollback(callback Function) {
	tx.steps = append(tx.steps, SimpleTransactionStep{
		RollbackFunc: callback,
	})
}

func (tx *Transaction) Execute() error {
	var i int
	var cause error
	for i, step := range tx.steps {
		cause = step.Execute()
		if cause != nil {
			logrus.Errorf("Encountered error on step %d of transaction: %s", i, cause.Error())
			break
		}
	}

	if cause != nil {
		logrus.Infof("Rolling back from step %d", i)
		for j := i; j >= 0; j-- {
			err := tx.steps[j].Rollback()
			if err != nil {
				logrus.Errorf("Ignoring error on step %d of rollback: %s", j, err.Error())
			}
		}
	}

	return cause
}
