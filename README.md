# 🇧🇷 Português

# Site de apresentação para Grupy-Sanca

## Alunos

- Eduardo Neves Gomes da Silva - 13822710
- Cauê Paiva Lira - 14675416
- Guilherme Antonio Costa Bandeira - 14575620

## Informações do projeto

Estamos desenvolvendo uma página web para o Grupy-Sanca, com o propósito de fortalecer sua presença online e reunir em um só lugar todas as informações relevantes do grupo. O site apresentará sua trajetória, redes sociais, fotos de eventos anteriores e anúncios de futuras atividades, além de oferecer integração com o sistema atual de inscrição em eventos.

### Objetivo
Fortalecer a presença online do Grupy-Sanca em um único hub simples e acolhedor. Apresentar o grupo, próximos eventos e a memória da comunidade de forma clara e convidativa.

### Público-alvo e jornada
Pessoas fora da bolha de tech e fora do público universitário, com perfil mais casual. Descobrem o site pela busca, entendem a história do grupo, conferem próximos eventos e se inscrevem.

### Páginas do site

1) **Início**
Recepção rápida com mensagem clara sobre a comunidade e botões para História e Eventos. Destaque para próximos eventos e links de redes sociais.

2) **História**
Resumo de como o Grupy-Sanca surgiu e seus principais marcos. Fotos históricas selecionadas para gerar identificação e confiança.

3) **Galeria**
Exploração visual de eventos passados em uma grade simples e intuitiva. Ao abrir uma foto, o visitante entende de qual evento ela veio e o contexto do registro.

4) **Código de Conduta**
Reforço dos valores de respeito, segurança e inclusão. Link direto para a versão oficial do Código de Conduta do Grupy-Sanca.

### Princípios de conteúdo e UX
Linguagem acessível e direta, foco na primeira visita, imagens otimizadas e navegação simples. Convites visíveis para participar de eventos e conhecer mais sobre o grupo.

## Arquitetura do Projeto

O projeto é dividido em duas partes principais:

### Front-end
- **Tecnologias:** React + TypeScript + Vite + Tailwind CSS + shadcn-ui
- **Deploy:** Hospedado na Vercel em [https://site-grupy.vercel.app/](https://site-grupy.vercel.app/)
- **Funcionalidades especiais:** Sistema de login para administradores com capacidades de edição de conteúdo e upload de mídia (imagens, eventos da timeline, textos)

### Back-end
- **Linguagem:** Go (Golang)
- **Arquitetura:** Clean Architecture (Arquitetura em Camadas) com Ports & Adapters
- **Estrutura:**
  - `/cmd` - Ponto de entrada da aplicação
  - `/internal/entities` - Entidades de domínio
  - `/internal/service` - Lógica de negócio e portas (interfaces)
  - `/internal/clients` - Adaptadores para serviços externos
  - `/internal/http/handlers` - Controllers HTTP
  - `/internal/repository` - Acesso a banco de dados
  - `/internal/gateway` - Gateway para object storage
- **Banco de dados:** Firebase Firestore
- **Storage:** Firebase Storage (para imagens)
- **API:** REST JSON em `/api/v1`

## Rodando o Projeto

### Front-end

#### Pré-requisitos
- Node.js 18+ e npm
- Git

#### Instalação
```sh
npm install
```

#### Execução do servidor local

```sh
npm run dev
```

Agora é só abrir http://127.0.0.1:5173 no seu browser

### Back-end

#### Pré-requisitos
- Go 1.22+
- Arquivo de credenciais do Firebase (JSON)

#### Configuração

1. Navegue até a pasta do backend:
```sh
cd backend
```

2. **IMPORTANTE:** Você precisa de um arquivo de credenciais do Firebase. Coloque o arquivo JSON de credenciais do Firebase Admin SDK na pasta `backend/` com o nome especificado nas configurações (ex: `sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json`).

3. Configure as variáveis de ambiente criando um arquivo `.env` na pasta `backend/`:
```sh
# Exemplo de .env
FIREBASE_CREDENTIALS_PATH=./sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
PORT=8080
```

#### Executando localmente

```sh
cd backend
go run cmd/server/main.go
```

O servidor estará disponível em `http://localhost:8080`

#### Endpoints principais da API

- `GET /api/v1/events` - Lista eventos
- `GET /api/v1/texts` - Lista textos
- `GET /api/v1/images` - Lista imagens
- `GET /api/v1/timelineentries` - Lista entradas da timeline
- `GET /api/v1/galery_events` - Lista eventos da galeria
- `GET /health` - Health check

## Deploy

### Front-end (Vercel)

O front-end está hospedado na Vercel e é atualizado automaticamente a cada push na branch `main`:
- **URL de produção:** [https://site-grupy.vercel.app/](https://site-grupy.vercel.app/)

### Back-end (Docker)

O back-end pode ser empacotado como uma imagem Docker e executado em qualquer servidor ou serviço de nuvem (GCP Cloud Run, AWS ECS, Azure Container Instances, etc.).

#### Build da imagem Docker

```sh
cd backend
docker build -t grupysanca-backend .
```

#### Executar o container localmente

```sh
docker run -p 8080:8080 \
  -e FIREBASE_CREDENTIALS_PATH=/app/credentials.json \
  -v /path/to/your/firebase-credentials.json:/app/credentials.json \
  grupysanca-backend
```

#### Deploy em Cloud Services

O container pode ser facilmente implantado em:
- **Google Cloud Run** (recomendado para integração com Firebase)
- **AWS ECS/Fargate**
- **Azure Container Instances**
- **Qualquer servidor com Docker**

## Tecnologias Usadas

### Front-end
- Vite
- TypeScript
- React
- React Router
- shadcn-ui
- Tailwind CSS
- Firebase Authentication

### Back-end
- Go (Golang)
- Firebase Firestore
- Firebase Storage
- Firebase Admin SDK

---

# 🇺🇸 English

# Grupy-Sanca Presentation Website

## Students

- Eduardo Neves Gomes da Silva - 13822710
- Cauê Paiva Lira - 14675416
- Guilherme Antonio Costa Bandeira - 14575620

## Project Information

We are developing a website for Grupy-Sanca, with the purpose of strengthening its online presence and bringing together in one place all relevant information about the group. The site will present its history, social media, photos from past events, and announcements of future activities, in addition to offering integration with the current event registration system.

### Objective
Strengthen Grupy-Sanca's online presence in a single simple and welcoming hub. Present the group, upcoming events, and the community's memory in a clear and inviting way.

### Target Audience and Journey
People outside the tech bubble and outside the university audience, with a more casual profile. They discover the site through search, understand the group's history, check out upcoming events, and sign up.

### Website Pages

1) **Home**
Quick welcome with a clear message about the community and buttons for History and Events. Highlight upcoming events and social media links.

2) **History**
Summary of how Grupy-Sanca started and its main milestones. Selected historical photos to generate identification and trust.

3) **Gallery**
Visual exploration of past events in a simple and intuitive grid. When opening a photo, the visitor understands which event it came from and the context of the record.

4) **Code of Conduct**
Reinforcement of values of respect, safety, and inclusion. Direct link to the official version of Grupy-Sanca's Code of Conduct.

### Content and UX Principles
Accessible and direct language, focus on first visit, optimized images, and simple navigation. Visible invitations to participate in events and learn more about the group.

## Project Architecture

The project is divided into two main parts:

### Front-end
- **Technologies:** React + TypeScript + Vite + Tailwind CSS + shadcn-ui
- **Deploy:** Hosted on Vercel at [https://site-grupy.vercel.app/](https://site-grupy.vercel.app/)
- **Special Features:** Admin login system with content editing and media upload capabilities (images, timeline events, texts)

### Back-end
- **Language:** Go (Golang)
- **Architecture:** Clean Architecture (Layered Architecture) with Ports & Adapters
- **Structure:**
  - `/cmd` - Application entry point
  - `/internal/entities` - Domain entities
  - `/internal/service` - Business logic and ports (interfaces)
  - `/internal/clients` - Adapters for external services
  - `/internal/http/handlers` - HTTP controllers
  - `/internal/repository` - Database access
  - `/internal/gateway` - Object storage gateway
- **Database:** Firebase Firestore
- **Storage:** Firebase Storage (for images)
- **API:** REST JSON at `/api/v1`

## Running the Project

### Front-end

#### Prerequisites
- Node.js 18+ and npm
- Git

#### Installation
```sh
npm install
```

#### Running the local server

```sh
npm run dev
```

Now just open http://127.0.0.1:5173 in your browser

### Back-end

#### Prerequisites
- Go 1.22+
- Firebase credentials file (JSON)

#### Configuration

1. Navigate to the backend folder:
```sh
cd backend
```

2. **IMPORTANT:** You need a Firebase credentials file. Place the Firebase Admin SDK credentials JSON file in the `backend/` folder with the name specified in the settings (e.g., `sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json`).

3. Configure environment variables by creating a `.env` file in the `backend/` folder:
```sh
# .env example
FIREBASE_CREDENTIALS_PATH=./sitegrupysanca-firebase-adminsdk-fbsvc-ff7567bd6e.json
PORT=8080
```

#### Running locally

```sh
cd backend
go run cmd/server/main.go
```

The server will be available at `http://localhost:8080`

#### Main API Endpoints

- `GET /api/v1/events` - List events
- `GET /api/v1/texts` - List texts
- `GET /api/v1/images` - List images
- `GET /api/v1/timelineentries` - List timeline entries
- `GET /api/v1/galery_events` - List gallery events
- `GET /health` - Health check

## Deployment

### Front-end (Vercel)

The front-end is hosted on Vercel and is automatically updated with each push to the `main` branch:
- **Production URL:** [https://site-grupy.vercel.app/](https://site-grupy.vercel.app/)

### Back-end (Docker)

The back-end can be packaged as a Docker image and run on any server or cloud service (GCP Cloud Run, AWS ECS, Azure Container Instances, etc.).

#### Building the Docker image

```sh
cd backend
docker build -t grupysanca-backend .
```

#### Running the container locally

```sh
docker run -p 8080:8080 \
  -e FIREBASE_CREDENTIALS_PATH=/app/credentials.json \
  -v /path/to/your/firebase-credentials.json:/app/credentials.json \
  grupysanca-backend
```

#### Cloud Services Deployment

The container can be easily deployed to:
- **Google Cloud Run** (recommended for Firebase integration)
- **AWS ECS/Fargate**
- **Azure Container Instances**
- **Any server with Docker**

## Technologies Used

### Front-end
- Vite
- TypeScript
- React
- React Router
- shadcn-ui
- Tailwind CSS
- Firebase Authentication

### Back-end
- Go (Golang)
- Firebase Firestore
- Firebase Storage
- Firebase Admin SDK
