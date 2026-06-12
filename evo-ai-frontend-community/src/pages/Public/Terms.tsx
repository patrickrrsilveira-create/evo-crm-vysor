import React from 'react';
import { Card } from '@evoapi/design-system';

const Terms: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <Card className="p-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Termos de Uso</h1>
          <div className="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300">
            <p className="mb-4">Última atualização: {new Date().toLocaleDateString('pt-BR')}</p>
            
            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">1. Aceitação dos Termos</h2>
            <p className="mb-4">
              Ao acessar e usar nossos serviços e nossa plataforma, operados pela VYSOR TECH LTDA,
              você concorda em cumprir e estar vinculado a estes Termos de Uso. Se você não concordar
              com qualquer parte destes termos, não deverá usar nossos serviços.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">2. Descrição do Serviço</h2>
            <p className="mb-4">
              A VYSOR TECH LTDA fornece um software de gestão de relacionamento com o cliente (CRM)
              que permite integrar múltiplos canais de comunicação, incluindo provedores externos como Meta (Facebook, Instagram, WhatsApp).
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">3. Contas de Usuário</h2>
            <p className="mb-4">
              Para usar nosso serviço, você deve criar uma conta. Você é responsável por manter a
              confidencialidade da sua conta e senha, e por restringir o acesso ao seu computador ou dispositivo.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">4. Integração com Plataformas de Terceiros</h2>
            <p className="mb-4">
              <ul className="list-disc pl-6 mt-2">
                <li>O serviço permite a integração de contas de terceiros (ex: páginas de Facebook, perfis de Instagram).</li>
                <li>O usuário deve cumprir integralmente os Termos de Serviço dessas plataformas (ex: Políticas de Privacidade e Termos da Meta).</li>
                <li>A VYSOR TECH LTDA não é responsável por eventuais bloqueios, suspensões ou falhas nos serviços destas plataformas externas.</li>
              </ul>
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">5. Uso Aceitável</h2>
            <p className="mb-4">
              Você concorda em usar nossos serviços apenas para fins lícitos. É proibido o uso da plataforma para
              enviar spam, conteúdo ilegal, ofensivo ou que viole direitos de terceiros. Reservamo-nos o direito de
              suspender ou encerrar contas que violem estas regras.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">6. Limitação de Responsabilidade</h2>
            <p className="mb-4">
              A VYSOR TECH LTDA não será responsável por quaisquer danos indiretos, incidentais, especiais ou
              consequentes decorrentes do uso ou da incapacidade de uso do serviço.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">7. Contato</h2>
            <p className="mb-4">
              Para qualquer dúvida sobre os Termos de Uso, entre em contato:
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

export default Terms;
