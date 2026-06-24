import React from 'react';
import { Card } from '@evoapi/design-system';

const DataDeletion: React.FC = () => {
  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900 py-12 px-4 sm:px-6 lg:px-8">
      <div className="max-w-3xl mx-auto">
        <Card className="p-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white mb-6">Exclusão de Dados</h1>
          <div className="prose dark:prose-invert max-w-none text-gray-700 dark:text-gray-300">
            <p className="mb-4">
              A VYSOR TECH LTDA respeita o seu direito à privacidade e o controle sobre os seus próprios dados. 
              De acordo com a Lei Geral de Proteção de Dados (LGPD) e as diretrizes de plataformas parceiras (como a Meta), você tem o direito de solicitar a exclusão de todos os seus dados pessoais do nosso sistema.
            </p>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">Como solicitar a exclusão dos seus dados</h2>
            <p className="mb-4">
              Se você deseja que todas as suas informações pessoais, histórico de conversas, nome, telefone ou e-mail sejam permanentemente removidos dos servidores do nosso CRM, siga os passos abaixo:
            </p>
            
            <ol className="list-decimal pl-6 mt-2 mb-6">
              <li className="mb-2">
                Envie um e-mail para <strong>contato@vysortech.com.br</strong> com o assunto "Solicitação de Exclusão de Dados".
              </li>
              <li className="mb-2">
                No corpo do e-mail, informe o seu Nome Completo e o Número de Telefone ou E-mail exato pelo qual você foi contatado ou entrou em contato conosco.
              </li>
              <li className="mb-2">
                Nossa equipe técnica irá processar a sua solicitação em até 48 horas úteis e todos os registros associados ao seu contato serão apagados de forma irreversível do banco de dados do Evolution CRM.
              </li>
            </ol>

            <h2 className="text-xl font-semibold mt-6 mb-3 text-gray-900 dark:text-white">Removendo o acesso do Aplicativo (Facebook/Instagram)</h2>
            <p className="mb-4">
              Se você é um usuário ou administrador e deseja remover as permissões do nosso aplicativo na sua conta do Facebook ou Instagram, você pode fazer isso diretamente nas configurações da sua rede social:
            </p>
            <ol className="list-decimal pl-6 mt-2 mb-6">
              <li>Acesse o seu Facebook e vá em <strong>Configurações e privacidade {'>'} Configurações</strong>.</li>
              <li>No menu lateral, clique em <strong>Segurança e login</strong> ou <strong>Integrações comerciais</strong>.</li>
              <li>Encontre o aplicativo "Vysor Tech" ou "Vysor CRM" na lista de aplicativos ativos.</li>
              <li>Clique em <strong>Remover</strong> para revogar permanentemente o acesso do nosso sistema à sua conta.</li>
            </ol>

            <p className="mt-8 text-sm text-gray-500">
              Última atualização: {new Date().toLocaleDateString('pt-BR')}
            </p>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default DataDeletion;
