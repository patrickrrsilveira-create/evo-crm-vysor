import React from 'react';
import { Card } from '@evoapi/design-system';

const Privacy: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <Card className="p-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Política de Privacidade</h1>
          <div className="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300">
            <p className="mb-4">Última atualização: {new Date().toLocaleDateString('pt-BR')}</p>
            
            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">1. Introdução</h2>
            <p className="mb-4">
              A VYSOR TECH LTDA ("nós", "nosso") respeita a sua privacidade e está comprometida em proteger
              seus dados pessoais. Esta política de privacidade informará como cuidamos dos seus dados pessoais
              quando você visita nosso site (vysortech.com.br) e utiliza nossos serviços.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">2. Dados que Coletamos</h2>
            <p className="mb-4">
              Podemos coletar, usar, armazenar e transferir diferentes tipos de dados pessoais sobre você, incluindo:
              <ul className="list-disc pl-6 mt-2">
                <li>Dados de Identidade (nome, sobrenome)</li>
                <li>Dados de Contato (endereço de e-mail, telefone)</li>
                <li>Dados Técnicos (endereço IP, tipo de navegador, sistema operacional)</li>
                <li>Dados de Uso (informações sobre como você usa nosso site e serviços)</li>
              </ul>
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">3. Como Usamos seus Dados</h2>
            <p className="mb-4">
              Usamos seus dados apenas quando a lei nos permite. Mais comumente, usaremos seus dados pessoais nas seguintes circunstâncias:
              <ul className="list-disc pl-6 mt-2">
                <li>Para executar o contrato que estamos prestes a celebrar ou celebramos com você.</li>
                <li>Quando for necessário para nossos interesses legítimos.</li>
                <li>Para cumprir uma obrigação legal ou regulatória.</li>
              </ul>
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">4. Integrações de Terceiros (Meta / Facebook / WhatsApp)</h2>
            <p className="mb-4">
              Nossos serviços oferecem integração com plataformas de terceiros, como o Facebook, Instagram e WhatsApp (operados pela Meta).
              Ao autorizar a conexão de nossos aplicativos à sua conta da Meta:
              <ul className="list-disc pl-6 mt-2">
                <li>Acessamos e armazenamos mensagens enviadas e recebidas nas páginas conectadas, com o propósito exclusivo de exibi-las em nosso CRM para atendimento.</li>
                <li>Não compartilhamos as mensagens ou dados de seus clientes com outros provedores externos que não sejam diretamente envolvidos no funcionamento do nosso CRM.</li>
                <li>Você pode revogar o acesso do nosso aplicativo a qualquer momento através das configurações da sua conta na respectiva plataforma terceira.</li>
              </ul>
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">5. Segurança de Dados</h2>
            <p className="mb-4">
              Implementamos medidas de segurança adequadas para evitar que seus dados pessoais sejam perdidos
              acidentalmente, usados ou acessados ​​de forma não autorizada, alterados ou divulgados.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">6. Contato</h2>
            <p className="mb-4">
              Se você tiver alguma dúvida sobre esta política de privacidade ou nossas práticas de privacidade, entre em contato conosco:
              <br/><br/>
              <strong>VYSOR TECH LTDA</strong><br/>
              E-mail: contato@vysortech.com.br<br/>
              Site: vysortech.com.br
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Privacy;
