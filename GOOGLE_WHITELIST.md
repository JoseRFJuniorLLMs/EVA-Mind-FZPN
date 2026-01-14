# Google Features Developer Whitelist

## Como habilitar Google Calendar para seu CPF:

1. Abra o arquivo `main.go`
2. Localize a variável `googleFeaturesWhitelist` (linha ~63)
3. Adicione seu CPF (somente números):

```go
googleFeaturesWhitelist = map[string]bool{
    "12345678900": true, // Exemplo
    "98765432100": true, // Seu CPF aqui
}
```

4. Reinicie o servidor

## Como funciona:

- Quando um idoso tenta usar `manage_calendar_event` (criar/listar eventos)
- O sistema verifica se o CPF dele está na whitelist
- Se **SIM** → Executa normalmente
- Se **NÃO** → Retorna: *"Google Calendar features are currently in beta and not available for your account."*

## Para liberar para todos (futuro):

Remova ou comente a verificação no `case "manage_calendar_event"`:

```go
// if !googleFeaturesWhitelist[client.CPF] {
//     return map[string]interface{}{...}
// }
```
