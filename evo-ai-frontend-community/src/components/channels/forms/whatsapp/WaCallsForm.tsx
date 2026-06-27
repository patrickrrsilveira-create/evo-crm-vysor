import React from 'react';
import { FormField } from '../../shared/FormField';
import { FormData } from '@/hooks/channels/useChannelForm';
import { sanitizeInboxName } from '@/utils/sanitizeName';
import { PhoneInput } from '@/components/shared/PhoneInput';

interface WaCallsFormProps {
  form: FormData;
  onFormChange: (key: string, value: string | boolean) => void;
}

export const WaCallsForm: React.FC<WaCallsFormProps> = ({ form, onFormChange }) => {
  const getStr = (key: string, fallback = ''): string =>
    typeof form[key] === 'string' ? (form[key] as string) : fallback;

  const handleDisplayNameChange = (value: string) => {
    onFormChange('display_name', value);
    onFormChange('name', sanitizeInboxName(value));
  };

  return (
    <div className="space-y-6">
      <FormField
        label="URL da API (WaCalls)"
        value={getStr('api_url', 'https://waha.vysortech.app.br')}
        onChange={(value: string) => onFormChange('api_url', value)}
        placeholder="https://waha.vysortech.app.br"
        type="url"
        required
      />

      <FormField
        label="Token de Acesso / API Key"
        value={getStr('api_key')}
        onChange={(value: string) => onFormChange('api_key', value)}
        placeholder="Insira a chave global da API"
        type="password"
        required
      />

      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <FormField
          label="Nome de Exibição"
          value={getStr('display_name')}
          onChange={handleDisplayNameChange}
          placeholder="Ex: Suporte WhatsApp"
          helpText="Nome amigável para identificação do canal"
          required
        />
        <FormField
          label="Nome do Canal"
          value={getStr('name')}
          onChange={(value: string) => onFormChange('name', value)}
          placeholder="Gerado automaticamente"
          helpText="Identificador único (letras minúsculas e sem espaços)"
          required
          readOnly
        />
      </div>

      <div>
        <label className="text-sm font-medium text-sidebar-foreground/80 block mb-1">
          Número do WhatsApp <span className="text-destructive">*</span>
        </label>
        <PhoneInput
          value={getStr('phone_number')}
          onChange={(value: string) => onFormChange('phone_number', value)}
          placeholder="+55 11 99999-9999"
          defaultCountry="BR"
        />
      </div>
    </div>
  );
};
