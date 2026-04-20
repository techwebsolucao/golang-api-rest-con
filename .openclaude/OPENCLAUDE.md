10 mandamentos de segurança e boas práticas Go (compacto)
1. Use go mod para dependências, evite GOPATH.
2. Não use panic para controle de fluxo, prefira erros.
3. Use context para timeout/cancelamento de goroutines.
4. Evite variáveis globais, passe dependências por função/struct.
5. Use defer para liberar recursos sempre.
6. Só use interface{} quando necessário.
7. Rode go vet e golint para checagem de código.
8. Prefira canais a sync.Mutex para concorrência.
9. Configure http.Client com timeout e cuidado.
10. Mantenha dependências sempre atualizadas.


