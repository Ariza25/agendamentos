export const API_ERROR_MESSAGES: Record<string, string> = {
  invalid_payload: 'Não foi possível ler os dados enviados.',
  invalid_input: 'Dados inválidos. Revise os campos e tente novamente.',
  not_found: 'Registro não encontrado.',
  internal_error: 'Não foi possível concluir a operação. Tente novamente.',
  unauthorized: 'Sessão expirada ou usuário não autenticado.',
  forbidden: 'Você não tem permissão para acessar este recurso.',
  time_conflict: 'Já existe um agendamento para este horário.',
  invalid_status_change: 'Transição de status inválida para este agendamento.',
  no_available_slots: 'Não há horários disponíveis para esta data.',
  booking_disabled: 'O agendamento online está desativado para este estabelecimento.',
  slug_already_used: 'Já existe um estabelecimento com este identificador de agenda.',
  slug_already_taken: 'Este link de agenda já está em uso por outro estabelecimento.',
  last_active_tenant: 'Não é possível remover o último estabelecimento ativo.',
  last_admin: 'Não é possível remover o último acesso administrativo.',
  email_already_registered: 'Já existe uma conta cadastrada com este e-mail.',
  invalid_credentials: 'E-mail ou senha incorretos.',
  invalid_reset_token: 'Link de redefinição inválido ou expirado.',
  invalid_reset_code: 'Código inválido ou expirado.',
  rate_limited: 'Muitas tentativas. Aguarde alguns minutos e tente novamente.',
  rate_limit_unavailable: 'Não foi possível validar a solicitação agora. Tente novamente.',
  auth_service_unavailable: 'Não foi possível processar a solicitação agora. Tente novamente.',
  mfa_required: 'Informe o código do aplicativo autenticador.',
  mfa_setup_failed: 'Não foi possível iniciar a configuração de dois fatores.',
  invalid_mfa_code: 'Código de verificação inválido.',
  invalid_mfa_confirmation: 'Senha ou código de verificação inválido.',
  session_list_failed: 'Não foi possível carregar os dispositivos conectados.',
  session_revoke_failed: 'Não foi possível encerrar a sessão.',
  invalid_profile: 'Não foi possível atualizar o perfil.',
  password_change_failed: 'Não foi possível alterar a senha.',
  service_unavailable: 'Serviço temporariamente indisponível.',
  appointment_done: 'Este agendamento já foi concluído e não pode mais ser alterado online.',
  cancellation_window: 'O prazo para cancelar ou reagendar este agendamento já expirou.',
  invalid_management_code: 'Código de agendamento inválido ou não encontrado.',
  notification_failed: 'Não foi possível entregar a notificação.',
  review_submitted: 'Esta avaliação já foi enviada.',
  signup_failed: 'Não foi possível concluir o cadastro. Verifique os dados e tente novamente.',
  connect_failed: 'Falha ao sincronizar o calendário.',
  invalid_rule: 'Regra de comissão inválida.',
  update_failed: 'Não foi possível atualizar as informações.',
  payout_failed: 'Não foi possível processar o repasse financeiro.',
  invalid_transaction: 'Transação financeira inválida.',
  invalid_product: 'Produto inválido.',
  invalid_id: 'Identificador obrigatório ou inválido.',
  invalid_movement: 'Movimentação de estoque inválida.',
  set_supplies_failed: 'Não foi possível vincular os insumos ao serviço.',
  missing_tenant: 'Estabelecimento não identificado.',
}

export const EXACT_MESSAGE_MAP: Record<string, string> = {
  'admin name is required': 'O seu nome completo é obrigatório.',
  'business name is required': 'O nome do estabelecimento é obrigatório.',
  'slug is required': 'O link da agenda é obrigatório.',
  'invalid slug': 'O link da agenda deve conter apenas letras minúsculas, números e traços.',
  'slug already taken': 'O link informado já está em uso por outro estabelecimento.',
  'email is required': 'O e-mail de acesso é obrigatório e deve ser válido.',
  'password must have at least 8 characters': 'A senha deve ter no mínimo 8 caracteres.',
  'Enter the full name.': 'Informe seu nome completo.',
  'Enter the email.': 'Informe seu e-mail.',
  'Enter a valid email.': 'Informe um e-mail válido.',
  'Enter a password.': 'Informe sua senha.',
  'An account with this email already exists.': 'Já existe uma conta cadastrada com este e-mail.',
  'Enter the current password.': 'Informe a senha atual.',
  'Enter the new password.': 'Informe a nova senha.',
  'Current password incorrect.': 'Senha atual incorreta.',
  'Invalid reset link.': 'Link de redefinição inválido.',
  'Enter your email.': 'Informe seu e-mail.',
  'Enter the 6-digit code.': 'Informe o código de 6 dígitos.',
  'Enter only the numbers of the code.': 'Informe apenas os números do código.',
  'Enter your password.': 'Informe sua senha.',
  'Incorrect email or password.': 'E-mail ou senha incorretos.',
  'Invalid verification code.': 'Código de verificação inválido.',
  'Session expired or user not authenticated.': 'Sessão expirada ou usuário não autenticado.',
  'You do not have permission to access this resource.': 'Você não tem permissão para acessar este recurso.',
  'Too many login attempts. Wait a few minutes.': 'Muitas tentativas de login. Aguarde alguns minutos.',
  'Too many attempts. Wait a few minutes and try again.': 'Muitas tentativas. Aguarde alguns minutos e tente novamente.',
  'Too many attempts. Please wait a moment and try again.': 'Muitas tentativas. Aguarde um momento e tente novamente.',
  'Internal server failure.': 'Ocorreu uma falha no servidor. Tente novamente em instantes.',
  'validation failed': 'Por favor, revise os dados informados.',
  'Failed to fetch': 'Falha na conexão com o servidor. Verifique sua internet.',
  'NetworkError when attempting to fetch resource.': 'Falha na conexão de rede. Verifique sua internet.',
  'The password must have at least 8 characters and contain uppercase letters, lowercase letters, numbers and a special character.':
    'A senha deve ter no mínimo 8 caracteres, contendo letras maiúsculas, minúsculas, números e caractere especial.',
}

const SUBSTRING_MESSAGE_PATTERNS: Array<[string, string]> = [
  ['admin name', 'O seu nome completo é obrigatório.'],
  ['business name', 'O nome do estabelecimento é obrigatório.'],
  ['slug already', 'Este link de agenda já está em uso por outro estabelecimento.'],
  ['invalid slug', 'O link da agenda deve conter apenas letras minúsculas, números e traços.'],
  ['email already exists', 'Já existe uma conta cadastrada com este e-mail.'],
  ['already exists', 'Já existe um registro com estes dados.'],
  ['password must have at least 8', 'A senha deve ter no mínimo 8 caracteres, contendo maiúsculas, minúsculas, números e caractere especial.'],
  ['incorrect email or password', 'E-mail ou senha incorretos.'],
  ['invalid credentials', 'E-mail ou senha incorretos.'],
  ['validation failed', 'Por favor, revise os campos destacados.'],
  ['create owner account', 'Não foi possível cadastrar o usuário administrador.'],
  ['create tenant', 'Não foi possível cadastrar o estabelecimento.'],
  ['rate limit', 'Muitas tentativas. Aguarde alguns minutos antes de tentar novamente.'],
  ['too many', 'Muitas tentativas. Aguarde alguns minutos antes de tentar novamente.'],
  ['failed to fetch', 'Falha na conexão com o servidor. Verifique sua internet.'],
  ['network', 'Falha de conexão com a rede. Verifique sua internet.'],
]

export function translateErrorMessage(rawMessage?: string | null): string {
  if (!rawMessage) {
    return 'Ocorreu um erro inesperado. Tente novamente.'
  }

  const trimmed = rawMessage.trim()
  if (EXACT_MESSAGE_MAP[trimmed]) {
    return EXACT_MESSAGE_MAP[trimmed]
  }

  const lower = trimmed.toLowerCase()
  for (const [pattern, translation] of SUBSTRING_MESSAGE_PATTERNS) {
    if (lower.includes(pattern)) {
      return translation
    }
  }

  return trimmed
}

export function getApiErrorMessage(code: string, fallback: string): string {
  if (API_ERROR_MESSAGES[code]) {
    return API_ERROR_MESSAGES[code]
  }
  return translateErrorMessage(fallback)
}

export function translateFieldErrors(fieldErrors?: Record<string, string>): Record<string, string> {
  if (!fieldErrors) {
    return {}
  }

  const translated: Record<string, string> = {}
  for (const [field, message] of Object.entries(fieldErrors)) {
    translated[field] = EXACT_MESSAGE_MAP[message] ?? translateErrorMessage(message)
  }
  return translated
}
