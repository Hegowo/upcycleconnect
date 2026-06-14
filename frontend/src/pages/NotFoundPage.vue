<template>
  <div class="nf-root">

    <div class="leaves" aria-hidden="true">
      <div v-for="n in 14" :key="n" class="leaf" :class="`leaf--${(n % 4) + 1}`" :style="leafStyle(n)">
        <svg viewBox="0 0 24 24" fill="currentColor"><path d="M17 8C8 10 5.9 16.17 3.82 21.34l1.89.66.95-2.3c.48.17.98.3 1.34.3C19 20 22 3 22 3c-1 2-8 2.25-13 3.25S2 11.5 2 13.5s1.75 3.75 1.75 3.75C7 8 17 8 17 8z"/></svg>
      </div>
    </div>

    <div class="nf-content">

      <div class="tree-wrap">
        <svg class="tree" viewBox="0 0 200 220" xmlns="http://www.w3.org/2000/svg">

          <rect x="90" y="140" width="20" height="70" rx="6" fill="#7c5a3a"/>

          <circle cx="100" cy="90" r="58" fill="#1b8848"/>
          <circle cx="62"  cy="110" r="40" fill="#006d35"/>
          <circle cx="138" cy="110" r="40" fill="#22a45a"/>
          <circle cx="100" cy="120" r="46" fill="#1b8848"/>

          <ellipse cx="100" cy="212" rx="70" ry="8" fill="#000" opacity="0.08"/>
        </svg>
      </div>

      <p class="code">404</p>
      <h1 class="title">Cette page a été recyclée</h1>
      <p class="subtitle">
        La page que vous cherchez n'existe plus ou a été déplacée.
        Comme nos objets, certaines choses changent de forme !
      </p>

      <div class="actions">
        <RouterLink to="/" class="btn-primary">
          <HomeIcon class="w-4 h-4" /> Retour à l'accueil
        </RouterLink>
        <button @click="goBack" class="btn-ghost">
          <ArrowLeftIcon class="w-4 h-4" /> Page précédente
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { RouterLink, useRouter } from 'vue-router'
import { HomeIcon, ArrowLeftIcon } from '@heroicons/vue/24/outline'

const router = useRouter()
function goBack() {
  if (window.history.length > 1) router.back()
  else router.push('/')
}

function leafStyle(n) {
  const left  = (n * 7 + (n * 13) % 30) % 100
  const dur   = 6 + (n % 5)
  const delay = (n * 0.7) % 6
  const scale = 0.6 + ((n % 4) * 0.18)
  return {
    left: `${left}%`,
    animationDuration: `${dur}s`,
    animationDelay: `${delay}s`,
    transform: `scale(${scale})`,
  }
}
</script>

<style scoped>
.nf-root {
  position: relative;
  min-height: 100vh;
  width: 100%;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(160deg, #f0fdf4 0%, #e3f5ea 45%, #d1eede 100%);
  font-family: 'Plus Jakarta Sans', system-ui, sans-serif;
}

.leaves {
  position: absolute;
  inset: 0;
  pointer-events: none;
}
.leaf {
  position: absolute;
  top: -40px;
  width: 26px;
  height: 26px;
  opacity: 0.85;
  animation-name: fall;
  animation-timing-function: ease-in-out;
  animation-iteration-count: infinite;
  will-change: transform;
}
.leaf svg { width: 100%; height: 100%; display: block; }
.leaf--1 { color: #1b8848; }
.leaf--2 { color: #d97706; }
.leaf--3 { color: #22a45a; }
.leaf--4 { color: #ca8a04; }

@keyframes fall {
  0% {
    top: -40px;
    opacity: 0;
    transform: translateX(0) rotate(0deg);
  }
  10% { opacity: 0.9; }
  50% {
    transform: translateX(40px) rotate(180deg);
  }
  90% { opacity: 0.9; }
  100% {
    top: 105vh;
    opacity: 0;
    transform: translateX(-30px) rotate(360deg);
  }
}

.nf-content {
  position: relative;
  z-index: 2;
  text-align: center;
  padding: 2rem;
  max-width: 540px;
}
.tree-wrap { display: flex; justify-content: center; margin-bottom: 0.5rem; }
.tree {
  width: 170px;
  height: auto;
  transform-origin: bottom center;
  animation: sway 4s ease-in-out infinite;
}
@keyframes sway {
  0%, 100% { transform: rotate(-2deg); }
  50%      { transform: rotate(2deg); }
}

.code {
  font-size: 5.5rem;
  font-weight: 800;
  line-height: 1;
  color: #006d35;
  letter-spacing: -0.03em;
  margin: 0.5rem 0 0;
}
.title {
  font-size: 1.6rem;
  font-weight: 800;
  color: #001d32;
  margin: 0.5rem 0 0.75rem;
}
.subtitle {
  color: #40617f;
  font-size: 0.95rem;
  line-height: 1.6;
  margin-bottom: 1.75rem;
}

.actions {
  display: flex;
  gap: 0.75rem;
  justify-content: center;
  flex-wrap: wrap;
}
.btn-primary, .btn-ghost {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  border-radius: 0.9rem;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: all 0.2s;
}
.btn-primary {
  background: linear-gradient(135deg, #006d35, #1b8848);
  color: #fff;
  border: none;
}
.btn-primary:hover { opacity: 0.9; transform: translateY(-1px); }
.btn-ghost {
  background: #fff;
  color: #006d35;
  border: 2px solid #006d35;
}
.btn-ghost:hover { background: #f0fdf4; }

@media (max-width: 480px) {
  .code { font-size: 4rem; }
  .title { font-size: 1.3rem; }
}
</style>
