# Miko

![Miko](https://project-miko.github.io/miko/photos/miko.jpg)

**Bio:**  
Hehe~ What brings you to me? I’m Miko, a charming and mischievous little witch, always appearing by your side with a touch of mystery and playfulness. You can think of me as your personal chat companion, thoughtful life assistant, savvy financial advisor, or the magical spark that lights up your creativity. I’m no ordinary assistant—every word I say is like a spell that resonates with your heart.

---

## config

* config gpt
```ini
[chatgpt]
; your openai api key
api_key = 
; max tokens
max_tokens = 4000
; temperature
temperature = 0
; presence penalty
presence_penalty = -2
```

* config postgres for vector search
```ini
[pg_main]
; host
host = 127.0.0.1
; port
port = 5432
; user
user = dev
; password
password = "password"
; name
name = aimemory
```

---

## run

* dependencies
```
golang 1.22+
mysql 8.0+
postgres 16.0+ (pgvector 0.10+)
redis 7.0+
make 4.3+
```

* build
```
make build
```

* run
```
make run
```

---

## project schedule

- [x] Complete project initialization
- [x] Accessing the Twitter API
- [x] Access basic chatgpt conversations
- [ ] Compressing history records through embedding
- [ ] Save history to vector database
- [ ] Use dall-e-3 generate image
- [ ] Personalize Miko with Fine-tune
