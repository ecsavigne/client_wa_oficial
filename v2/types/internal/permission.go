package internal

// map[Key]map[Label, Description]value
type PERMISSION_TYPE map[string]map[string]string

// Get returns the value of the permission map for the given key and subkey(Label, Description).
// If the key or subkey does not exist, it returns an empty string.
func (self PERMISSION_TYPE) Get(key, subkey string) string {
	return self[key][subkey]
}

var permissions PERMISSION_TYPE = PERMISSION_TYPE{
	"ads_management": {
		"Label":       "Gerenciamento de Anúncios",
		"Description": "Permite ler e administrar contas publicitárias para criar campanhas, gerenciar anúncios e obter métricas de anúncios.",
	},
	"ads_read": {
		"Label":       "Leitura de Anúncios",
		"Description": "Permite acessar a API de estatísticas de anúncios para obter relatórios de anúncios e enviar eventos web ao Facebook.",
	},
	"attribution_read": {
		"Label":       "Leitura de Atribuição",
		"Description": "Permite acessar a API de atribuição para obter dados de relatórios de atribuição para análise e painéis personalizados.",
	},
	"business_management": {
		"Label":       "Gerenciamento de Negócios",
		"Description": "Permite ler e escrever com a API do administrador de negócios para gerenciar ativos comerciais como contas publicitárias.",
	},
	"catalog_management": {
		"Label":       "Gerenciamento de Catálogo",
		"Description": "Permite criar, ler, atualizar e eliminar catálogos de produtos para desenvolver soluções comerciais e de gerenciamento de inventário.",
	},
	"commerce_account_manage_orders": {
		"Label":       "Gerenciar Pedidos da Conta de Comércio",
		"Description": "Permite ler e atualizar pedidos da conta de comércio, acessar notificações de webhook e gerenciar pedidos em nome de clientes.",
	},
	"commerce_account_read_orders": {
		"Label":       "Ler Pedidos da Conta de Comércio",
		"Description": "Permite ler pedidos da conta de comércio e usar endereços de email para fins de marketing com consentimento.",
	},
	"commerce_account_read_reports": {
		"Label":       "Ler Relatórios da Conta de Comércio",
		"Description": "Permite ler dados de relatórios financeiros para gerar relatórios personalizados de impostos, conciliação de caixa e reembolsos.",
	},
	"commerce_account_read_settings": {
		"Label":       "Ler Configurações da Conta de Comércio",
		"Description": "Permite ler configurações da conta de comércio, como canais conectados, opções de envio e locais de processamento.",
	},
	"commerce_manage_accounts": {
		"Label":       "Gerenciar Contas de Comércio",
		"Description": "Permite criar e administrar contas de comércio, ativar novos canais de vendas e associar aplicativos com contas de comércio.",
	},
	"facebook_branded_content_ads_brand": {
		"Label":       "Anúncios de Conteúdo de Marca Facebook",
		"Description": "Permite ler publicações do Facebook onde a conta está etiquetada como parceiro de pagamento e gerenciar permissões para publicar anúncios de parceria.",
	},
	"facebook_creator_marketplace_discovery": {
		"Label":       "Descoberta Marketplace de Criadores Facebook",
		"Description": "Permite descobrir criadores de conteúdo e acessar dados de estatísticas para avaliação em campanhas de marca e reconhecimento de criadores.",
	},
	"email": {
		"Label":       "Email",
		"Description": "Permite ler o endereço de email principal de uma pessoa para comunicação e autenticação na aplicação.",
	},
	"gaming_user_locale": {
		"Label":       "Idioma do Usuário Gaming",
		"Description": "Permite conhecer o idioma de preferência de um usuário ao jogar no Facebook para exibir a interface do jogo no idioma preferido.",
	},
	"instagram_basic": {
		"Label":       "Instagram Básico",
		"Description": "Permite ler informações e conteúdo multimídia do perfil de uma conta de Instagram para obter metadados básicos.",
	},
	"instagram_branded_content_ads_brand": {
		"Label":       "Anúncios de Conteúdo de Marca Instagram",
		"Description": "Permite ler publicações onde a conta está etiquetada como parceiro de pagamento e gerenciar permissões para publicar anúncios de parceria.",
	},
	"instagram_branded_content_brand": {
		"Label":       "Conteúdo de Marca Instagram",
		"Description": "Permite agregar, eliminar e visualizar criadores na lista de criadores aprovados para uma marca específica.",
	},
	"instagram_branded_content_creator": {
		"Label":       "Conteúdo de Marca Criador Instagram",
		"Description": "Permite ler e modificar o estado de promoção de conteúdo específico de um criador em Instagram.",
	},
	"instagram_business_basic": {
		"Label":       "Instagram Business Básico",
		"Description": "Permite ler informações e conteúdo multimídia do perfil de uma conta profissional do Instagram.",
	},
	"instagram_business_content_publish": {
		"Label":       "Publicação de Conteúdo Business Instagram",
		"Description": "Permite criar publicações orgânicas com fotos e vídeos no feed em nome de uma conta profissional.",
	},
	"instagram_business_manage_comments": {
		"Label":       "Gerenciar Comentários Business Instagram",
		"Description": "Permite criar, eliminar e ocultar comentários em uma conta profissional do Instagram.",
	},
	"instagram_business_manage_messages": {
		"Label":       "Gerenciar Mensagens Business Instagram",
		"Description": "Permite acessar mensagens em uma conta profissional do Instagram e usar ferramentas de CRM externas.",
	},
	"instagram_content_publish": {
		"Label":       "Publicação de Conteúdo Instagram",
		"Description": "Permite criar publicações orgânicas com fotos e vídeos no feed em nome de um usuário comercial.",
	},
	"instagram_creator_marketplace_discovery": {
		"Label":       "Descoberta Marketplace de Criadores Instagram",
		"Description": "Permite descobrir criadores e acessar dados de estatísticas como apresentação, contagem de seguidores e alcance.",
	},
	"instagram_creator_marketplace_messaging": {
		"Label":       "Mensageria Marketplace de Criadores Instagram",
		"Description": "Permite obter conversas de parceria e ID de mensageria para enviar mensagens de parceria pagas a criadores.",
	},
	"instagram_manage_comments": {
		"Label":       "Gerenciar Comentários Instagram",
		"Description": "Permite criar, eliminar, ocultar comentários e responder a menções em contas vinculadas a páginas.",
	},
	"instagram_manage_contents": {
		"Label":       "Gerenciar Conteúdos Instagram",
		"Description": "Permite eliminar publicações, histórias e reels em nome de uma conta vinculada a uma página.",
	},
	"instagram_manage_events": {
		"Label":       "Gerenciar Eventos Instagram",
		"Description": "Permite registrar eventos em contas do Instagram e enviar dados de atividade ao Meta para segmentação e otimização.",
	},
	"instagram_manage_insights": {
		"Label":       "Gerenciar Insights Instagram",
		"Description": "Permite acessar estatísticas de contas profissionais do Instagram, incluindo metadados, dados e insights de histórias.",
	},
	"instagram_manage_messages": {
		"Label":       "Gerenciar Mensagens Instagram",
		"Description": "Permite ler e responder mensagens diretas do Instagram para gerenciar conversas com clientes.",
	},
	"instagram_shopping_tag_products": {
		"Label":       "Etiquetar Produtos Shopping Instagram",
		"Description": "Permite etiquetar conteúdo multimídia com etiquetas de produtos, gerenciar etiquetas e apelar rejeições.",
	},
	"instagram_manage_upcoming_events": {
		"Label":       "Gerenciar Eventos Próximos Instagram",
		"Description": "Permite ler, criar e atualizar eventos próximos em contas do Instagram administradas pelo usuário.",
	},
	"leads_retrieval": {
		"Label":       "Recuperação de Clientes Potenciais",
		"Description": "Permite recuperar e ler informações registradas através de formulários de anúncios para contatar pessoas interessadas ou para plataformas de CRM extrair dados em nome dos anunciantes.",
	},
	"manage_app_solutions": {
		"Label":       "Gerenciar Soluções de Aplicativo",
		"Description": "Permite obter uma lista de aplicativos que um usuário pode gerenciar e fazer chamadas de API em nome desses aplicativos para criar e administrar soluções de parceiros.",
	},
	"manage_fundraisers": {
		"Label":       "Gerenciar Campanhas de Arrecadação",
		"Description": "Permite criar, atualizar e ler uma campanha de arrecadação de fundos e suas doações em nome de um usuário.",
	},
	"marketing_messages_messenger": {
		"Label":       "Mensagens de Marketing Messenger",
		"Description": "Permite criar, gerenciar e enviar mensagens de marketing pagas no Messenger para pessoas que optaram por receber anúncios e promoções.",
	},
	"pages_events": {
		"Label":       "Eventos de Página",
		"Description": "Permite registrar eventos em nome de páginas do Facebook e enviar dados para segmentação de anúncios, otimização e relatórios.",
	},
	"pages_manage_ads": {
		"Label":       "Gerenciar Anúncios de Página",
		"Description": "Permite criar e gerenciar anúncios associados à página, incluindo anúncios de cliques em plataformas de mensageria comercial.",
	},
	"pages_manage_cta": {
		"Label":       "Gerenciar CTA de Página",
		"Description": "Permite executar funções de publicação e exclusão para administrar botões de chamada à ação em uma página do Facebook.",
	},
	"pages_manage_instant_articles": {
		"Label":       "Gerenciar Artigos Instantâneos",
		"Description": "Permite criar e atualizar artigos instantâneos em nome das páginas do Facebook que o usuário administra.",
	},
	"pages_manage_engagement": {
		"Label":       "Gerenciar Engajamento de Página",
		"Description": "Permite criar, editar ou eliminar comentários na página para ajudar a gerenciar e moderar conteúdo.",
	},
	"pages_manage_metadata": {
		"Label":       "Gerenciar Metadados de Página",
		"Description": "Permite se inscrever em webhooks de atividade da página, recebê-los e atualizar configurações da página.",
	},
	"pages_manage_posts": {
		"Label":       "Gerenciar Publicações de Página",
		"Description": "Permite criar, editar ou eliminar publicações, fotos e vídeos de uma página.",
	},
	"pages_messaging": {
		"Label":       "Mensageria de Página",
		"Description": "Permite acessar e gerenciar conversas da página no Messenger para criar experiências interativas e enviar mensagens de atendimento.",
	},
	"pages_read_engagement": {
		"Label":       "Ler Engajamento de Página",
		"Description": "Permite ler conteúdo publicado pela página, acessar dados de seguidores e metadados para administrar a página.",
	},
	"pages_read_user_content": {
		"Label":       "Ler Conteúdo do Usuário na Página",
		"Description": "Permite ler conteúdo gerado por usuários e eliminar comentários para gerenciar a página.",
	},
	"pages_show_list": {
		"Label":       "Mostrar Lista de Páginas",
		"Description": "Permite acessar a lista de páginas que um usuário administra e verificar se é o proprietário.",
	},
	"pages_user_gender": {
		"Label":       "Gênero do Usuário na Página",
		"Description": "Permite acessar o gênero de um usuário através de uma página para personalizar experiências e usar pronomes corretos.",
	},
	"pages_user_locale": {
		"Label":       "Idioma do Usuário na Página",
		"Description": "Permite acessar o idioma de um usuário através de uma página para personalizar conteúdo e respostas.",
	},
	"pages_user_timezone": {
		"Label":       "Zona Horária do Usuário na Página",
		"Description": "Permite acessar a zona horária de um usuário através de uma página para enviar mensagens e conteúdo em horários apropriados.",
	},
	"pages_utility_messaging": {
		"Label":       "Mensagens de Utilidade de Página",
		"Description": "Permite acessar e enviar modelos de mensagens de utilidade de uma página pelo Messenger.",
	},
	"private_computation_access": {
		"Label":       "Acesso à Computação Privada",
		"Description": "Permite acessar ambientes de computação privada do Meta para realizar processamento seguro de dados e medições com proteção de privacidade.",
	},
	"public_profile": {
		"Label":       "Perfil Público",
		"Description": "Permite ler os campos do perfil público padrão para autenticar usuários e oferecer experiências personalizadas.",
	},
	"publish_video": {
		"Label":       "Publicar Vídeo",
		"Description": "Permite publicar vídeos em direto na biografia, grupos, eventos ou página do usuário da aplicação.",
	},
	"read_audience_network_insights": {
		"Label":       "Ler Insights da Audience Network",
		"Description": "Permite acessar dados de insights do Audience Network e integrar informações de desempenho em painéis de análise.",
	},
	"read_insights": {
		"Label":       "Ler Insights",
		"Description": "Permite ler dados de estatísticas de páginas, apps e domínios web para integrar com ferramentas de análise próprias.",
	},
	"read_page_mailboxes": {
		"Label":       "Leitura de Caixas de Entrada da Página",
		"Description": "Permite ler mensagens e conversas recebidas na caixa de entrada de Páginas do Facebook para exibição e gerenciamento no sistema.",
	},
	"threads_basic": {
		"Label":       "Threads Básico",
		"Description": "Permite obter informações do perfil de Threads de um usuário e o conteúdo de mídia e texto postado em Threads.",
	},
	"threads_business_basic": {
		"Label":       "Threads Business Básico",
		"Description": "Permite obter o ID da conta de Threads associado a uma conta do Instagram para criar anúncios em Threads.",
	},
	"threads_content_publish": {
		"Label":       "Publicação de Conteúdo Threads",
		"Description": "Permite criar e publicar conteúdo em nome de um perfil de Threads para gerenciar a presença do usuário.",
	},
	"threads_delete": {
		"Label":       "Exclusão Threads",
		"Description": "Permite eliminar as publicações de Threads de um usuário para gerenciar sua presença na plataforma.",
	},
	"threads_keyword_search": {
		"Label":       "Busca por Palavra-chave Threads",
		"Description": "Permite buscar e recuperar conteúdo com palavras-chave específicas em Threads e publicar respostas.",
	},
	"threads_location_tagging": {
		"Label":       "Marcação de Localização Threads",
		"Description": "Permite buscar localizações públicas e publicar mídia com localização marcada em Threads.",
	},
	"threads_manage_insights": {
		"Label":       "Gerenciar Insights Threads",
		"Description": "Permite obter acesso a estatísticas de um perfil de Threads e métricas de postagens individuais.",
	},
	"threads_manage_mentions": {
		"Label":       "Gerenciar Menções Threads",
		"Description": "Permite buscar conteúdo onde o usuário é mencionado em Threads e gerenciar sua presença social.",
	},
	"threads_manage_replies": {
		"Label":       "Gerenciar Respostas Threads",
		"Description": "Permite criar respostas, ocultar ou mostrar respostas a um thread e controlar quem pode responder.",
	},
	"threads_profile_discovery": {
		"Label":       "Descoberta de Perfil Threads",
		"Description": "Permite acessar perfis de contas públicas de Threads e postagens públicas para análises competitivas.",
	},
	"threads_read_replies": {
		"Label":       "Ler Respostas Threads",
		"Description": "Permite ler respostas aos threads de um usuário para gerenciar sua presença em Threads.",
	},
	"user_birthday": {
		"Label":       "Data de Nascimento do Usuário",
		"Description": "Permite ler a data de nascimento de uma pessoa do seu perfil do Facebook para fornecer conteúdo relevante à idade.",
	},
	"user_friends": {
		"Label":       "Amigos do Usuário",
		"Description": "Permite obter uma lista de amigos de uma pessoa que usam a aplicação para personalizar a experiência.",
	},
	"user_gender": {
		"Label":       "Gênero do Usuário",
		"Description": "Permite ler o gênero de uma pessoa do seu perfil do Facebook para personalizar a experiência e interpretar pronomes.",
	},
	"user_gender_range": {
		"Label":       "Faixa Etária do Usuário",
		"Description": "Permite acessar a faixa etária de uma pessoa do seu perfil do Facebook para aplicar restrições de idade ao conteúdo.",
	},
	"user_hometown": {
		"Label":       "Cidade de Origem do Usuário",
		"Description": "Permite ler a cidade de origem de uma pessoa do seu perfil do Facebook para oferecer uma experiência personalizada.",
	},
	"user_likes": {
		"Label":       "Afinidades do Usuário",
		"Description": "Permite ler uma lista de páginas do Facebook que o usuário marcou como gostando para personalizar a experiência e análise de preferências.",
	},
	"user_link": {
		"Label":       "Link do Perfil do Usuário",
		"Description": "Permite acessar a URL do perfil do Facebook de uma pessoa para que os usuários possam visitar perfis de outros usuários.",
	},
	"user_location": {
		"Label":       "Localização do Usuário",
		"Description": "Permite ler a cidade do campo de localização do perfil do Facebook de uma pessoa para oferecer uma experiência personalizada.",
	},
	"user_messenger_contact": {
		"Label":       "Contato via Messenger",
		"Description": "Permite que uma empresa entre em contato com uma pessoa via Messenger para enviar mensagens iniciais, atualizações pós-compra e atualizações de conta.",
	},
	"user_photos": {
		"Label":       "Fotos do Usuário",
		"Description": "Permite ler as fotos que uma pessoa enviou para o Facebook para criar álbuns, compartilhar e editar fotos.",
	},
	"user_posts": {
		"Label":       "Publicações do Usuário",
		"Description": "Permite acessar as publicações que um usuário fez em seu perfil para criar álbuns, compartilhar memórias e analisar conteúdo.",
	},
	"user_videos": {
		"Label":       "Vídeos do Usuário",
		"Description": "Permite ler uma lista de vídeos que uma pessoa enviou para personalizar a exibição e permitir edição de conteúdo de vídeo.",
	},
	"whatsapp_business_management": {
		"Label":       "Gerenciador WhatsApp Business",
		"Description": "Permite gerenciar ativos do WhatsApp Business e exibir análises da conta no portal do cliente.",
	},
	"whatsapp_business_manage_events": {
		"Label":       "Gerenciador de Eventos WhatsApp Business",
		"Description": "Permite registrar eventos em contas do WhatsApp Business e enviar dados de atividade ao Meta para segmentação, otimização e relatórios de anúncios.",
	},
	"whatsapp_business_messaging": {
		"Label":       "Mensageria WhatsApp Business",
		"Description": "Permite enviar mensagens e fazer chamadas de WhatsApp, enviar e recuperar conteúdo multimídia, gerenciar informações do perfil de negócios e registrar números de telefone com o Meta.",
	},
}

func GetPermissionType() PERMISSION_TYPE {
	return permissions
}
