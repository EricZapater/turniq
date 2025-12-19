# Annex estètic — Criteris UI Frontend Turniq

## Objectiu visual

Dotar el frontend de **Turniq** d’una estètica:

- Neta
- Professional
- Industrial / B2B
- Discreta, sense artificis visuals innecessaris

Ha de transmetre:

- Fiabilitat
- Simplicitat
- Eina de treball, no producte consumer

---

## Principis de disseny

1. Claredat abans que decoració
2. Poques decisions visuals
3. Jerarquia clara (què és important i què no)
4. Zero soroll visual
5. Consistència absoluta

---

## Paleta de colors

Utilitzar una paleta reduïda:

### Color principal

- Blau fosc o gris industrial
- Ús: botons primaris, focus, accions principals

### Colors neutres

- Blanc o gris molt clar per fons
- Gris mitjà per textos secundaris

### Color d’error

- Vermell sobri (no saturat)
- Ús exclusiu per errors i alerts

**Restriccions**:

- No utilitzar gradients
- No utilitzar colors cridaners
- No utilitzar ombres agressives

---

## Tipografia

- Utilitzar la tipografia per defecte de PrimeVue o una sans-serif neutra
- Jerarquia clara de mides:

  - Títol
  - Subtítol
  - Text normal

- No utilitzar més de dos pesos de font

---

## Layout del login

### Estructura

- Pantalla centrada verticalment i horitzontalment
- Targeta (Card) amb:

  - Títol: "Accés"
  - Formulari
  - Botó principal clar

### Espaiat

- Padding generós
- Camps ben separats
- Aire visual

---

## Components PrimeVue autoritzats

- Card: contenidor del login
- InputText: email
- Password: contrasenya
- Button: acció principal
- Toast: errors i feedback

**Restriccions**:

- No utilitzar components experimentals
- No utilitzar animacions complexes

---

## Botons

- Un únic botó primari
- Estat disabled mentre s’envia el formulari
- Feedback clar en hover i focus

---

## Missatges d’error

### Errors de validació

- Discrets
- Sota el camp corresponent

### Errors globals

- Sempre via Toast
- Text curt, directe, no tècnic

---

## Comportament visual

- Mostrar loader o disabled state durant el login
- Evitar salts de layout
- No mostrar informació tècnica (status codes, stack traces)

---

## Restriccions estètiques

- No aspecte "startup trendy"
- No animacions flashy
- No iconografia innecessària
- No fons foscos

---

## Criteri final d’acceptació visual

Si el login sembla:

- una eina interna d’empresa industrial,
- usable per personal no tècnic,
- clara a la primera ullada,

**l’estètica és correcta.**
