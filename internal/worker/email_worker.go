package worker

import (
	"log"
	"nineshop-be/pkg/email"
)

type BackgroundWorker interface {
	EnqueueEmailTask(toEmail, otpCode string)
}

type EmailTask struct {
	ToEmail string
	OTPCode string
}

type EmailWorkerPool struct {
	taskQueue    chan EmailTask
	emailService email.EmailService
	workerCount  int
}

func NewEmailWorkerPool(bufferSize int, workerCount int, emailSvc email.EmailService) BackgroundWorker {
	return &EmailWorkerPool{
		taskQueue:    make(chan EmailTask, bufferSize),
		emailService: emailSvc,
		workerCount:  workerCount,
	}
}

func (w *EmailWorkerPool) Start() {
	for i := 0; i < w.workerCount; i++ {
		go func(workerID int) {
			for task := range w.taskQueue {
				log.Printf("[EmailWorker #%d] Đang gửi OTP tới: %s", workerID, task.ToEmail)

				err := w.emailService.SendOTPEmail(task.ToEmail, task.OTPCode)
				if err != nil {
					log.Printf("[EmailWorker #%d] Lỗi khi gửi email tới %s: %v", workerID, task.ToEmail, err)
				} else {
					log.Printf("[EmailWorker #%d] Gửi email thành công tới %s", workerID, task.ToEmail)
				}
			}
		}(i)
	}
}

func (w *EmailWorkerPool) EnqueueEmailTask(toEmail, otpCode string) {
	w.taskQueue <- EmailTask{
		ToEmail: toEmail,
		OTPCode: otpCode,
	}
}
