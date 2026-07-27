---
icon: lucide/palette
---

# Design Language

Meshbox Studio is an archive. Most of what fills the screen is other people's
work: renders, photographs of finished prints, cover images pulled out of
Printables exports. The interface's job is to hold that content and get out of
its way.

Everything below follows from that one idea.

## Principles

**The content is the only thing allowed to be colourful.** Chrome is neutral.
Prints photograph warm — orange PLA, red PETG, wood-fill browns — so the
interface accent is cool. Warm hues on a Meshbox screen mean *heat*, and heat
means something is happening right now.

**Black is a real colour, not "very dark grey".** The dark theme canvas is
`#000000`. On an OLED panel those pixels are switched off: no light, no power,
and a thumbnail sitting on top of them has nothing to compete with.

**Depth comes from light, never from shadow.** A drop shadow against `#000` is
invisible. Surfaces separate by getting *lighter* and by carrying a translucent
white hairline.

**Badges mark exceptions.** A status chip on every card teaches the eye to skip
status chips. Only unusual states get a hue.

**Measured values are monospaced.** Print times, filament weights, layer
heights, file sizes, nozzle diameters, version strings. They line up in columns
and they don't reflow when a digit changes.

---

## Colour

### Roles

| Role | Family | Means | Appears as |
| --- | --- | --- | --- |
| `primary` | cyan | "you can act here" | Solid buttons, links, focus rings, active nav |
| `printing` | amber | Work happening right now | Badge + breathing dot |
| `success` | emerald | A print completed | Outcome badges |
| `warning` | orange | Attention, nothing broken yet | Alert bars, storage meter ≥75% |
| `error` | rose | Failed, or about to destroy something | Failed prints, hover state of destructive buttons |
| `info` | blue | Neutral information | Queued iterations |
| `neutral` | graphite | Everything else | All chrome |

`primary` and `printing` are the two that matter, and they are deliberately
opposites: **cool means you can do something, warm means the machine is doing
something.**

`printing` is a first-class Nuxt UI colour, registered in
[`nuxt.config.ts`](https://github.com/meshbox-studio/meshbox-studio/blob/main/frontend/nuxt.config.ts),
so `<UBadge color="printing">` and `text-printing` both work. Nothing except
in-progress work may use it.

!!! note "Why `warning` and `printing` are both warm"
    Orange and amber sit close together on purpose — both mean "look here". They
    never share a component, and form keeps them apart: `warning` only ever
    appears as an alert bar or a meter fill, `printing` only ever as a badge with
    a breathing dot.

### Graphite

The neutral ramp is a true neutral pulled a few degrees warm, so it sits under
an amber badge without turning green and never competes with the cool accent.

| | Hex | | Hex |
| --- | --- | --- | --- |
| `graphite-50` | `#fafaf8` | `graphite-500` | `#737269` |
| `graphite-100` | `#f2f2ef` | `graphite-600` | `#57564f` |
| `graphite-200` | `#e6e5e0` | `graphite-700` | `#403f3a` |
| `graphite-300` | `#d1d0ca` | `graphite-800` | `#292824` |
| `graphite-400` | `#a09f97` | `graphite-900` | `#191917` |
| | | `graphite-950` | `#0c0c0b` |

### The dark surface ladder

Five steps, each one a decision about how far something has lifted off the
canvas. Reach for the lowest one that works.

| Token | Value | Used for |
| --- | --- | --- |
| `--mb-canvas` | `#000000` | The page itself, the sidebar, the navbar |
| `--ui-bg` | `#0c0c0b` | Cards, popovers, dropdowns, modals |
| `--ui-bg-muted` | `#141413` | Input wells, recessed rows, thumbnail backing |
| `--ui-bg-elevated` | `#1a1a18` | Hover states, soft badges |
| `--ui-bg-accented` | `#262622` | Pressed states, chips, meter tracks |

The sidebar and navbar take **no background of their own**. They sit on the
canvas and are separated by a hairline. A panel that paints itself grey is what
makes an OLED theme look like an ordinary dark theme.

### Hairlines

Borders in dark mode are translucent white, not a fixed grey, so a single token
reads correctly over every surface in the ladder above.

```css
--ui-border:          color-mix(in oklab, #fff  9%, transparent);
--ui-border-muted:    color-mix(in oklab, #fff  6%, transparent);
--ui-border-accented: color-mix(in oklab, #fff 15%, transparent);
```

### Text, and why nothing is pure white

Pure white on pure black blooms on OLED and smears during scroll. Headings stop
at `#f2f2ef`, which still clears 18:1.

| Token | Dark | Contrast on `#000` | Used for |
| --- | --- | --- | --- |
| `--ui-text-highlighted` | `#f2f2ef` | 18.4:1 | Headings, project titles |
| `--ui-text` | `#d1d0ca` | 14.2:1 | Body copy |
| `--ui-text-toned` | `#a09f97` | 8.9:1 | Secondary copy |
| `--ui-text-muted` | `#8a8982` | 6.0:1 | Metadata, labels |
| `--ui-text-dimmed` | `#7b7a73` | 4.9:1 | Timestamps, counts |

Every step clears WCAG AA for body text, including the dimmest one. There is no
tier below `dimmed`; if text needs to be quieter than 4.9:1, it should not be on
screen.

In light mode `primary` resolves to `cyan-700` rather than `cyan-500` — the
lighter shade cannot carry white label text at 4.5:1.

---

## Type

**Geist** for the interface, **Geist Mono** for measured values. Both are
self-hosted variable fonts in `frontend/public/fonts/geist/`.

Anything that is a measurement gets `.mb-measure`, which switches to mono and
turns on tabular figures and slashed zero:

```html
<span class="mb-measure">4h 12m</span>
<span class="mb-measure">0.2 mm</span>
<span class="mb-measure">327.0 GB / 1006.9 GB</span>
```

`<code>`, `<kbd>` and `<samp>` get this automatically.

Headings are `font-semibold tracking-tight`. Nothing above `text-xl` in the
dashboard — this is a workspace, not a landing page.

---

## Shape, spacing, motion

`--ui-radius: 0.5rem` drives everything: `xs` 4px, `sm` 8px, `md` 12px, `lg`
16px, `xl` 24px. Cards land at 16px, buttons at 12px.

Spacing is Tailwind's 4px scale. Card padding `p-4 sm:p-6`, grid gaps
`gap-4 sm:gap-6`, related items `gap-1.5` to `gap-3`.

Motion is short and only ever confirms something happened: 150ms for hover and
colour, 500ms for a meter filling. The single looping animation in the system is
`mb-live`, a 2.4s breath on the printing indicator — noticeable in peripheral
vision without demanding attention. `prefers-reduced-motion` disables all of it
globally.

---

## Domain utilities

### `.plate`

A build-plate grid, used behind every thumbnail, 3D viewport, and empty media
slot. It reads as a printer bed, and it gives transparent PNG renders something
to sit on instead of dissolving into the canvas.

```html
<img class="plate size-16 rounded-md border border-default object-cover">
```

### `.mb-live`

The breathing opacity loop. Applied to the dot inside a `printing` badge, and to
anything else that is genuinely in progress.

---

## Component rules

**One primary action per view.** It is solid `primary`, and in the dashboard it
is the rightmost thing in the navbar. If a second solid accent button appears on
a screen, one of them is wrong.

**Destructive actions are quiet at rest.** `color="neutral" variant="ghost"` with
`hover:text-error`. Red at rest trains people to ignore red. Only the
confirmation step is allowed to be solid `error`.

**State badges go through `<ProjectStateBadge>`**, never a hand-rolled `<UBadge>`.
It is the only place the state→colour→variant mapping lives, alongside
`frontend/app/utils/project.ts`.

Hue carries attention; **form carries completeness**. `Draft` and `Active` share
the neutral hue and are told apart by variant — `outline` reads as unfinished,
`subtle` reads as settled.

**Thumbnails** are `object-cover` on `.plate`, with `rounded-md` and a hairline.
Never a shadow.

**Empty states** are a `UAlert` with `variant="soft" color="neutral"` and a verb:
say what to do next, not that there is nothing here.

---

## Extending it

The whole system lives in two files:

- `frontend/app/assets/css/main.css` — tokens, ramp, fonts, utilities
- `frontend/app/app.config.ts` — role→family mapping

Before adding a colour, check whether an existing role already means what you
need. The palette is small on purpose: seven roles, one neutral ramp, and a rule
that warm means heat. Most new UI needs a new *form*, not a new hue.
