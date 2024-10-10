# Desafio Full Cycle: Client-Server-API

Olá dev, tudo bem?

Neste desafio vamos aplicar o que aprendemos sobre webserver http, contextos,
banco de dados e manipulação de arquivos com Go.

Você precisará nos entregar dois sistemas em Go:

- client.go

- server.go

Os requisitos para cumprir este desafio são:

O client.go deverá realizar uma requisição HTTP no server.go solicitando a cotação do dólar ✔.

O server.go deverá consumir a API contendo o câmbio de Dólar e Real no endereço: https://economia.awesomeapi.com.br/json/last/USD-BRL  ✔ e em seguida deverá retornar no formato JSON o resultado para o cliente ✔.

Usando o package "context" ✔, o server.go deverá registrar no banco de dados SQLite cada cotação recebida ✔, sendo que o timeout máximo para chamar a API de cotação do dólar deverá ser de 200ms ✔ e o timeout máximo para conseguir persistir os dados no banco deverá ser de 10ms ✔.

O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON) ✔. Utilizando o package "context" ✔, o client.go terá um timeout máximo de 300ms para receber o resultado do server.go ✔.

Os 3 contextos deverão retornar erro nos logs caso o tempo de execução seja insuficiente ✔.

O client.go terá que salvar a cotação atual em um arquivo "cotacao.txt" no formato: Dólar: {valor} ✔

O endpoint necessário gerado pelo server.go para este desafio será: /cotacao ✔ e a porta a ser utilizada pelo servidor HTTP será a 8080 ✔.

Ao finalizar, envie o link do repositório para correção.

Obs: os arquivos gerados `cotacao.txt` e `cotacao.db` não foram incluidos no repositório.

- resposta da `awesomeapi`

```json
{
  "USDBRL": {
    "code": "USD",
    "codein": "BRL",
    "name": "Dólar Americano/Real Brasileiro",
    "high": "5.604",
    "low": "5.5306",
    "varBid": "-0.0224",
    "pctChange": "-0.4",
    "bid": "5.5712",
    "ask": "5.5722",
    "timestamp": "1728565880",
    "create_date": "2024-10-10 10:11:20"
  }
}
```

- saida `server.go`

```sh
antonio@DG15:~/DEV/full-cycle/client-server-api/server$ go run server.go 
2024/10/10 10:07:12 received: 
{"code":"USD","codein":"BRL","name":"Dólar Americano/Real Brasileiro","high":"5.604","low":"5.5306","varBid":"-0.02","pctChange":"-0.36","bid":"5.5736","ask":"5.5746","timestamp":"1728565603","create_date":"2024-10-10 10:06:43"}
2024/10/10 10:07:12 sent: 
{"bid":"5.5736"}
^C2024/10/10 10:07:23 server: shutting down
antonio@DG15:~/DEV/full-cycle/client-server-api/server$ 
```

- saída `cliente.go`

```sh
antonio@DG15:~/DEV/full-cycle/client-server-api/client$ go run client.go 
2024/10/10 10:07:12 recebido: 
{"bid":"5.5736"}
antonio@DG15:~/DEV/full-cycle/client-server-api/client$ 
```

- `cotacao.txt`

```txt
 Dólar: {5.5736}
```
