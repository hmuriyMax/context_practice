package userip

import (
	"context"
	"fmt"
	"net"
	"net/http"
)

// Тип ключа не экспортируется
// для предотвращения конфликтов
// с ключами контекста, определенными в других пакетах.
type key int

// userIPkey - это контекстный ключ
// для IP-адреса пользователя.
// Его нулевое значение произвольно.
// Если этот пакет определит другие контекстные ключи,
// они будут иметь разные целочисленные значения.
const userIPKey key = 0

func FromRequest(req *http.Request) (net.IP, error) {
	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("userip: %q is not IP:port", req.RemoteAddr)
	}
	return net.IP(ip), nil
}

func NewContext(ctx context.Context, userIP net.IP) context.Context {
	return context.WithValue(ctx, userIPKey, userIP)
}

func FromContext(ctx context.Context) (net.IP, bool) {
	// ctx.Value возвращает nil,
	// если ctx не имеет значения для ключа;
	// утверждение типа net.IP
	// возвращает ok=false для nil.
	userIP, ok := ctx.Value(userIPKey).(net.IP)
	return userIP, ok
}
