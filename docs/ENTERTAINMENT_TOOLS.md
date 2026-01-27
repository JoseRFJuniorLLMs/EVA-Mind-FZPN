# Ferramentas de Entretenimento - EVA-Mind-FZPN

**Documento:** ENTERTAINMENT-001
**Versão:** 1.0
**Data:** 2026-01-27
**Autor:** Jose R F Junior

---

## Visão Geral

30 ferramentas de entretenimento e bem-estar projetadas especificamente para idosos. Todas as ferramentas são ativadas por voz ou texto e adaptadas às capacidades do paciente.

---

## Categorias

| Categoria | Quantidade | Descrição |
|-----------|------------|-----------|
| Música e Áudio | 6 | Músicas, rádio, sons relaxantes |
| Jogos Cognitivos | 6 | Quiz, memória, exercícios mentais |
| Histórias | 5 | Narrativas, reminiscência, notícias |
| Bem-estar | 6 | Meditação, respiração, exercícios |
| Social/Família | 4 | Mensagens, fotos, conexão |
| Utilidades | 3 | Tempo, receitas, diário |

---

## 1. Música e Áudio (6 ferramentas)

### 1.1 play_nostalgic_music
**Tocar Música Nostálgica**

Toca músicas da época do paciente (anos 50-80).

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `decade` | string | 1950s, 1960s, 1970s, 1980s, any |
| `artist` | string | Nome do artista |
| `genre` | string | mpb, samba, bossa_nova, sertanejo, forro, bolero |
| `mood` | string | alegre, calma, romantica, animada, nostalgica |

**Exemplos de ativação:**
- "coloca uma música"
- "quero ouvir Roberto Carlos"
- "toca algo dos anos 60"
- "música pra relaxar"

---

### 1.2 play_radio_station
**Sintonizar Rádio**

Sintoniza estações de rádio online.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `station_type` | string | news, music, religious, local, sports |
| `station_name` | string | Nome da estação (CBN, Jovem Pan, etc) |

**Exemplos:**
- "liga a rádio"
- "quero ouvir notícias"
- "coloca a CBN"

---

### 1.3 nature_sounds
**Sons da Natureza**

Reproduz sons relaxantes para meditação e sono.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `sound_type` | string | rain, ocean, forest, birds, fireplace, river, thunderstorm, wind |
| `duration_minutes` | integer | Duração (padrão: 30) |
| `volume` | string | low, medium, high |

**Exemplos:**
- "som de chuva"
- "barulho do mar"
- "sons pra dormir"

---

### 1.4 audiobook_reader
**Ler Audiobook**

Lê livros em voz alta com controle de velocidade.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | play, pause, resume, stop, list, search |
| `book_title` | string | Título do livro |
| `chapter` | integer | Número do capítulo |
| `speed` | string | slow, normal, fast |

---

### 1.5 podcast_player
**Tocar Podcast**

Reproduz podcasts por categoria.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | play, pause, resume, list, search |
| `category` | string | health, history, humor, spirituality, news, culture |

---

### 1.6 religious_content
**Conteúdo Religioso**

Orações, terços, reflexões bíblicas, hinos.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `content_type` | string | prayer, rosary, bible_reflection, hymn, mass, meditation |
| `religion` | string | catholic, evangelical, spiritist, generic |
| `specific_prayer` | string | Oração específica (Pai Nosso, Ave Maria, Salmo 23) |

**Exemplos:**
- "reza comigo"
- "quero ouvir o terço"
- "lê um salmo"

---

## 2. Jogos Cognitivos (6 ferramentas)

### 2.1 play_trivia_game
**Quiz de Conhecimentos**

Perguntas e respostas sobre diversos temas.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start, answer, hint, skip, score, end |
| `theme` | string | brazil_history, music, geography, culture, sports, random |
| `difficulty` | string | easy, medium, hard |
| `answer` | string | Resposta do paciente |

---

### 2.2 memory_game
**Jogo da Memória**

Exercício de memória por voz - EVA diz sequência, paciente repete.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start, repeat, check, score, end |
| `game_type` | string | numbers, words, colors, objects |

---

### 2.3 word_association
**Associação de Palavras**

Jogo de associação - estimula conexões neurais.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start, respond, end |
| `category` | string | general, food, places, people, objects |

---

### 2.4 brain_training
**Treino Cerebral**

Exercícios cognitivos variados.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `exercise_type` | string | math, sequences, categories, opposites, analogies |
| `difficulty` | string | very_easy, easy, medium |
| `action` | string | start, answer, hint, next, end |

---

### 2.5 complete_the_lyrics
**Complete a Letra**

Jogo musical - EVA canta, paciente completa.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start, answer, skip, hint, score |
| `decade` | string | 1950s, 1960s, 1970s, 1980s, mixed |

---

### 2.6 riddles_and_jokes
**Charadas e Piadas**

Humor leve e adequado para idosos.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `content_type` | string | joke, riddle, tongue_twister, funny_story |
| `theme` | string | general, animals, family, daily_life, classic |

---

## 3. Histórias e Narrativas (5 ferramentas)

### 3.1 story_generator
**Gerador de Histórias**

Cria histórias personalizadas com nomes de familiares.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `story_type` | string | adventure, romance, family, childhood, fantasy, historical |
| `include_family` | boolean | Incluir nomes de familiares |
| `length` | string | short, medium, long |
| `setting` | string | Cenário (fazenda, cidade, praia) |

---

### 3.2 reminiscence_therapy
**Terapia de Reminiscência**

Conversa terapêutica sobre memórias do passado.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `theme` | string | childhood, youth, marriage, career, travels, holidays, family, friends |
| `use_music` | boolean | Usar músicas como gatilho |
| `use_photos` | boolean | Usar fotos do paciente |

**Benefícios clinicamente validados:**
- Melhora humor e bem-estar
- Reduz sintomas depressivos
- Fortalece identidade
- Estimula memória de longo prazo

---

### 3.3 biography_writer
**Escritor de Biografia**

Ajuda a construir biografia como legado para família.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start_session, continue, read_back, export, add_photo |
| `life_chapter` | string | birth_childhood, youth, love_marriage, career, parenthood, wisdom, legacy |

---

### 3.4 read_newspaper
**Ler Notícias**

Lê manchetes evitando conteúdo negativo.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `category` | string | general, sports, entertainment, health, local, positive |
| `source` | string | g1, uol, folha, estadao, local |
| `detail_level` | string | headlines, summary, full |

---

### 3.5 daily_horoscope
**Horóscopo do Dia**

Mensagens sempre positivas e motivacionais.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `sign` | string | aries, taurus, gemini, cancer, leo, virgo, libra, scorpio, sagittarius, capricorn, aquarius, pisces |

---

## 4. Bem-estar e Saúde (6 ferramentas)

### 4.1 guided_meditation
**Meditação Guiada**

Diferentes técnicas de meditação.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `technique` | string | mindfulness, body_scan, visualization, gratitude, loving_kindness, sleep |
| `duration_minutes` | integer | 5, 10, 15, 20 |
| `background_sound` | string | none, nature, music, bells |

---

### 4.2 breathing_exercises
**Exercícios de Respiração**

Técnicas validadas para ansiedade.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `technique` | string | 4-7-8, box_breathing, diaphragmatic, calming, energizing |
| `cycles` | integer | Número de ciclos (padrão: 5) |

**Técnica 4-7-8:**
1. Inspire por 4 segundos
2. Segure por 7 segundos
3. Expire por 8 segundos

---

### 4.3 chair_exercises
**Exercícios na Cadeira**

Exercícios físicos leves sentado.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `body_part` | string | full_body, arms, legs, neck, back, hands |
| `duration_minutes` | integer | Duração |
| `intensity` | string | gentle, moderate |

---

### 4.4 sleep_stories
**Histórias para Dormir**

Narrativas calmas para induzir sono.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `story_theme` | string | nature, journey, countryside, ocean, garden, clouds |
| `include_breathing` | boolean | Incluir pausas para respiração |

---

### 4.5 gratitude_journal
**Diário de Gratidão**

Prática diária de gratidão - 3 coisas boas.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | add_entry, read_today, read_week, read_random |
| `gratitude_items` | string | Coisas pelas quais está grato |

---

### 4.6 motivational_quotes
**Frases Motivacionais**

Citações de grandes pensadores.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `theme` | string | general, courage, love, faith, wisdom, perseverance, happiness |
| `author_type` | string | any, saints, philosophers, writers, brazilian |

---

## 5. Social e Família (4 ferramentas)

### 5.1 voice_capsule
**Cápsula de Voz**

Grava mensagens para família.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | record, play_back, send_now, schedule, list |
| `recipient` | string | Nome do familiar |
| `scheduled_date` | string | Data (YYYY-MM-DD) |
| `occasion` | string | birthday, holiday, just_because, anniversary, encouragement |

---

### 5.2 birthday_reminder
**Lembrete de Aniversários**

Gerencia datas de aniversário.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | check_today, check_week, check_month, add, list_all |
| `person_name` | string | Nome da pessoa |
| `date` | string | Data (DD/MM) |

---

### 5.3 family_tree_explorer
**Explorar Árvore Genealógica**

Navega pela genealogia e conta histórias.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | explore, add_person, add_story, view_tree, find_relation |
| `person_name` | string | Nome da pessoa |
| `relation` | string | Relação (avó materna, tio paterno) |

---

### 5.4 photo_slideshow
**Apresentação de Fotos**

Mostra fotos antigas com narração.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | start, pause, next, previous, stop, comment |
| `album` | string | childhood, wedding, family, travels, career, recent, all |
| `with_music` | boolean | Música de fundo |

---

## 6. Utilidades Diárias (3 ferramentas)

### 6.1 weather_chat
**Conversa sobre o Tempo**

Previsão do tempo com dicas.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `location` | string | Cidade |
| `forecast_type` | string | now, today, tomorrow, week |

---

### 6.2 cooking_recipes
**Receitas Culinárias**

Receitas brasileiras passo a passo.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | search, start_recipe, next_step, repeat_step, list_ingredients |
| `dish_type` | string | main, dessert, soup, salad, drink, snack |
| `difficulty` | string | easy, medium |

---

### 6.3 voice_diary
**Diário de Voz**

Grava pensamentos e reflexões.

| Parâmetro | Tipo | Opções |
|-----------|------|--------|
| `action` | string | record, play_today, play_date, play_random, list_recent |
| `date` | string | Data (YYYY-MM-DD) |
| `tag` | string | thought, memory, dream, gratitude, worry, plan |

---

## Permissões por Persona

| Ferramenta | Companion | Clinical | Emergency | Educator |
|------------|-----------|----------|-----------|----------|
| Música | ✅ | ❌ | ❌ | ❌ |
| Jogos | ✅ | ❌ | ❌ | ✅ |
| Histórias | ✅ | ❌ | ❌ | ✅ |
| Bem-estar | ✅ | ✅ (limitado) | ❌ | ✅ |
| Social | ✅ | ❌ | ❌ | ❌ |
| Utilidades | ✅ | ❌ | ❌ | ❌ |

---

## Arquivos Relacionados

| Arquivo | Descrição |
|---------|-----------|
| `migrations/017_entertainment_tools_seed.sql` | Seed de todas as ferramentas |
| `internal/tools/handlers.go` | Implementação dos handlers |
| `internal/tools/definitions.go` | Definições Go das ferramentas |

---

## Aprovações

| Função | Nome | Data |
|--------|------|------|
| Criador/Admin | Jose R F Junior | 2026-01-27 |

---

**Documento controlado - Versão 1.0**
