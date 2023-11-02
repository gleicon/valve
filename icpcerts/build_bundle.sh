#!/bin/sh

mkdir tmp
wget http://acraiz.icpbrasil.gov.br/credenciadas/CertificadosAC-ICP-Brasil/ACcompactado.zip
unzip  ACcompactado.zip -d tmp/
cat /etc/ssl/certs/icp-brasil/{AC_Certisign_RFB_G5.crt,AC_Secretaria_da_Receita_Federal_do_Brasil_v4.crt,ICP-Brasilv5.crt} > chain.pem
