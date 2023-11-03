![Valve](images/valve-image.jpg)

### Valve

Proxy compativel com mTLS para uso com certificados digitais compativeis com ICP Brasil (eCPF e eCNPJ)

Dados do certificado digitalOs dados extraidos de eCPF e eCNPJ, precisam da cadeia de certificado do [ICP Brasil](https://www.gov.br/iti/pt-br/assuntos/repositorio).

Estes dados são passados para o upstream (aplicação) por um header http -X-CERTIFICATE-CN, desta forma os serviços não precisam acessar o certificado. Se não houver certificado presente, o header X-CERTIFICATE-DETECTED: off será passado.

O diretório icpcerts tem um script e instruções para criar e atualizar os arquivos necessários a partir da fonte original.

### Build

$ make

### Run

As opções disponívels são:
-a endereço do proxy (IP)
-p porta do proxy
-u upstream server (pode usar um httpbin ou request catcher para testar)

Certificado do servidor - precisa ser vállido para o dominio que vai responder. 
-c seu certificado (pode ser a do letsencrypt)
-k sua chave privada (pode ser a do letsencrypt)

CACerts do Brasil do ICP (pode servir para outros casos semelhantes como proxies privados com CA local)
-a caminho para seu CACerts

Por exemplo usei meu certificado para https://ctofieldguide.com e coloquei a linha "127.0.0.1 ctofieldguide.com" no /etc/hosts. Executei o proxy com:

$ valve -c mycerts/cert13.pem -k mycerts/privkey13.pem -u https://nheco.requestcatcher.com/

E abri https://ctofieldguide.com. No request catcher deve ter um request com um header a mais contendo o valor do CN do certificado do cliente.


### Diagrama funcional

![arch](images/valve_arch.png)


