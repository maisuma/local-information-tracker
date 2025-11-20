package watcher

// Notifier は通知を管理する構造体
type Notifier struct {
	notifyChan chan string
}

// NewNotifier は新しい Notifier を作成する
func NewNotifier(bufferSize int) *Notifier { // コンストラクタ
	return &Notifier{
		notifyChan: make(chan string, bufferSize),
	}
}

// Notify は通知を送信する
func (n *Notifier) Notify(message string) {
	//stringで送っているが、必要に応じて構造体に変更可能
	n.notifyChan <- message
}

// StartListening は通知を受信して処理を実行する
// func (n *Notifier) StartListening(ctx context.Context, handler func(string)) {
// 	go func() {
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case msg := <-n.notifyChan:
// 				handler(msg)
// 			}
// 		}
// 	}()
// }
