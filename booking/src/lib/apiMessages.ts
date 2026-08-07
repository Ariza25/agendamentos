export const API_ERROR_MESSAGES: Record<string, string> = {
  invalid_payload: 'Não foi possível ler os dados enviados.',
  invalid_input: 'Dados inválidos. Revise os campos e tente novamente.',
  not_found: 'Registro não encontrado.',
  internal_error: 'Não foi possível concluir a operação. Tente novamente.',
  unauthorized: 'Sessão expirada ou usuário não autenticado.',
  forbidden: 'Você não tem permissão para acessar este recurso.',
  time_conflict: 'Este horário acabou de ser ocupado. Escolha outro horário disponível.',
  invalid_status_change: 'Transição de status inválida para este agendamento.',
  no_available_slots: 'Não há horários disponíveis para esta data.',
  booking_disabled: 'O agendamento online está desativado para este estabelecimento.',
  slug_already_used: 'Já existe um estabelecimento com este identificador.',
  last_active_tenant: 'Não é possível remover o último estabelecimento ativo.',
  last_admin: 'Não é possível remover o último acesso administrativo.',
  email_already_registered: 'Já existe uma conta cadastrada com este email.',
  invalid_credentials: 'Email ou senha incorretos.',
  invalid_reset_token: 'Link de redefinição inválido ou expirado.',
  invalid_reset_code: 'Código inválido ou expirado.',
  rate_limited: 'Muitas tentativas. Aguarde alguns minutos e tente novamente.',
  rate_limit_unavailable: 'Não foi possível validar a solicitação agora. Tente novamente.',
  auth_service_unavailable: 'Não foi possível processar a solicitação agora. Tente novamente.',
  mfa_required: 'Informe o código do aplicativo autenticador.',
  mfa_setup_failed: 'Não foi possível iniciar a configuração.',
  invalid_mfa_code: 'Código de verificação inválido.',
  invalid_mfa_confirmation: 'Senha ou código de verificação inválido.',
  session_list_failed: 'Não foi possível carregar os dispositivos conectados.',
  session_revoke_failed: 'Não foi possível encerrar a sessão.',
  invalid_profile: 'Não foi possível atualizar o perfil.',
  password_change_failed: 'Não foi possível alterar a senha.',
  service_unavailable: 'Serviço temporariamente indisponível.',
  appointment_done: 'Este agendamento não pode mais ser gerenciado online.',
  cancellation_window: 'O prazo para cancelar ou reagendar este agendamento já passou.',
  invalid_management_code: 'Código de agendamento inválido ou não encontrado.',
  notification_failed: 'Não foi possível entregar a notificação.',
  review_submitted: 'Esta avaliação já foi enviada.',
  customer_not_found: 'Não encontramos agendamentos para este telefone.',
}

export const FIELD_ERROR_MESSAGES: Record<string, string> = {
  'Enter the current password.': 'Informe a senha atual.',
  'Enter the new password.': 'Informe a nova senha.',
  'Current password incorrect.': 'Senha atual incorreta.',
  'Invalid reset link.': 'Link de redefinição inválido.',
  'Enter your email.': 'Informe seu email.',
  'Enter the 6-digit code.': 'Informe o código de 6 dígitos.',
  'Enter only the numbers of the code.': 'Informe apenas os números do código.',
  'Enter your password.': 'Informe sua senha.',
  'Incorrect email or password.': 'Email ou senha incorretos.',
  'Invalid verification code.': 'Código de verificação inválido.',
  'Enter the full name.': 'Informe o nome completo.',
  'Enter the email.': 'Informe o email.',
  'Enter a valid email.': 'Informe um email válido.',
  'Enter a password.': 'Informe uma senha.',
  'An account with this email already exists.': 'Já existe uma conta cadastrada com este email.',
  'The password must have at least 8 characters and contain uppercase letters, lowercase letters, numbers and a special character.': 'A senha deve ter no mínimo 8 caracteres, além de conter: letras maiúsculas, minúsculas, números e caractere especial.',
}

export function getApiErrorMessage(code: string, fallback: string): string {
  return API_ERROR_MESSAGES[code] ?? fallback
}

export function translateFieldErrors(fieldErrors?: Record<string, string>): Record<string, string> {
  if (!fieldErrors) {
    return {}
  }

  const translated: Record<string, string> = {}
  for (const [field, message] of Object.entries(fieldErrors)) {
    translated[field] = FIELD_ERROR_MESSAGES[message] ?? message
  }
  return translated
}
