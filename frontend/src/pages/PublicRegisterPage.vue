<template>
  <div class="min-h-screen flex flex-col bg-[#f7f9ff]">

    <header class="px-6 py-4 flex items-center justify-between max-w-[1280px] mx-auto w-full">
      <RouterLink to="/" class="flex items-center">
        <img src="/logoentier.png" class="h-10 w-auto" :alt="t('public.layout.logoAlt')" />
      </RouterLink>
      <button class="text-[#40617f] text-base font-medium px-4 py-2 rounded-xl hover:bg-[#edf4ff] transition">
        {{ t('public.register.headerHelp') }}
      </button>
    </header>

    <main class="flex-1 flex items-center justify-center px-4 py-12">
      <div class="w-full max-w-[896px] flex flex-col gap-12">

        <div class="flex items-center justify-center gap-3">
          <div
            v-for="i in (accountType === 'provider' ? 3 : 2)"
            :key="i"
            class="h-2 w-12 rounded-full transition-all duration-300"
            :style="i - 1 < stepIndex
              ? 'background:#006d35'
              : i - 1 === stepIndex
                ? 'background:#006d35'
                : 'background:#d8eaff'"
          />
        </div>

        <Transition name="step" mode="out-in">

          <div v-if="step === 'type'" key="type" class="flex flex-col gap-12">
            <div class="flex flex-col items-center gap-4 text-center">
              <h1 class="font-jakarta font-extrabold text-[#001d32] text-[30px] md:text-[48px] leading-[1.1] tracking-[-0.025em]">
                {{ t('public.register.step1Title') }}
              </h1>
              <p class="text-[#40617f] text-lg leading-[28px] max-w-[512px]">
                {{ t('public.register.step1Subtitle') }}
              </p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 pt-4">

              <button
                @click="accountType = 'user'"
                class="relative flex flex-col justify-between items-start p-10 rounded-[24px] border-2 text-left transition-all duration-200 hover:shadow-xl bg-white shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)]"
                :style="accountType === 'user'
                  ? 'border-color: rgba(0,109,53,0.4);'
                  : 'border-color: transparent;'"
              >
                <Transition name="badge">
                  <div v-if="accountType === 'user'" class="absolute top-6 right-6 w-8 h-8 bg-[#006d35] rounded-full flex items-center justify-center">
                    <CheckIcon class="w-3.5 h-3.5 text-white" />
                  </div>
                </Transition>

                <div class="w-16 h-16 rounded-[24px] bg-[#d8eaff] flex items-center justify-center mb-6">
                  <UserIcon class="w-5 h-5 text-[#40617f]" />
                </div>

                <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-3">{{ t('public.register.typeIndividualTitle') }}</h3>

                <p class="text-[#40617f] text-base leading-[26px] mb-8">
                  {{ t('public.register.typeIndividualDesc') }}
                </p>

                <div class="flex items-center gap-2 text-[#006d35] font-semibold text-base">
                  {{ t('public.register.btnSelect') }}
                  <ChevronRightIcon class="w-3 h-3" />
                </div>
              </button>

              <button
                @click="accountType = 'provider'"
                class="relative flex flex-col justify-between items-start p-10 rounded-[24px] border-2 text-left transition-all duration-200 hover:shadow-xl bg-white shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)]"
                :style="accountType === 'provider'
                  ? 'border-color: rgba(0,109,53,0.4);'
                  : 'border-color: transparent;'"
              >
                <Transition name="badge">
                  <div v-if="accountType === 'provider'" class="absolute top-6 right-6 w-8 h-8 bg-[#006d35] rounded-full flex items-center justify-center">
                    <CheckIcon class="w-3.5 h-3.5 text-white" />
                  </div>
                </Transition>

                <div class="w-16 h-16 rounded-[24px] bg-[#1b8848] flex items-center justify-center mb-6">
                  <WrenchScrewdriverIcon class="w-6 h-6 text-white" />
                </div>

                <h3 class="font-jakarta font-bold text-[#001d32] text-2xl mb-3">{{ t('public.register.typeProviderTitle') }}</h3>

                <p class="text-[#40617f] text-base leading-[26px] mb-8">
                  {{ t('public.register.typeProviderDesc') }}
                </p>

                <div class="flex items-center gap-2 text-[#006d35] font-semibold text-base">
                  {{ accountType === 'provider' ? t('public.register.btnSelected') : t('public.register.btnSelect') }}
                  <ChevronRightIcon v-if="accountType !== 'provider'" class="w-3 h-3" />
                </div>
              </button>
            </div>

            <div class="flex flex-col items-center gap-6 pt-4">
              <button
                @click="goToNext"
                :disabled="!accountType"
                class="w-80 py-5 rounded-xl font-jakarta font-bold text-lg transition-all duration-200"
                :style="accountType
                  ? 'background:#006d35;color:white;cursor:pointer;'
                  : 'background:#e2e8f0;color:#94a3b8;cursor:not-allowed;'"
              >
                {{ t('public.register.btnContinue') }}
              </button>
              <div class="flex items-center gap-2 text-[#40617f] text-sm">
                <LockClosedIcon class="w-3.5 h-3.5" />
                {{ t('public.register.dataProtected') }}
              </div>
            </div>
          </div>

          <div v-else-if="step === 'siret'" key="siret" class="max-w-lg mx-auto w-full">
            <button @click="step = 'type'" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#001d32] transition mb-6">
              <ArrowLeftIcon class="w-4 h-4" /> {{ t('public.register.back') }}
            </button>

            <div class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] p-10">
              <div class="flex items-center gap-4 mb-8">
                <div class="w-12 h-12 rounded-xl bg-[#edf4ff] flex items-center justify-center">
                  <MagnifyingGlassIcon class="w-5 h-5 text-[#40617f]" />
                </div>
                <div>
                  <h2 class="font-jakarta font-bold text-[#001d32] text-xl">{{ t('public.register.siretSearchTitle') }}</h2>
                  <p class="text-[#40617f] text-sm">{{ t('public.register.siretSearchSubtitle') }}</p>
                </div>
              </div>

              <div class="mb-4">
                <label class="block text-sm font-medium text-[#001d32] mb-2">
                  {{ t('public.register.siretLabel') }} <span class="text-xs font-normal text-[#40617f]">{{ t('public.register.siretLabelHint') }}</span>
                </label>
                <div class="relative">
                  <input
                    v-model="siretInput"
                    @input="onSiretInput"
                    type="text"
                    maxlength="14"
                    :placeholder="t('public.register.siretPlaceholder')"
                    class="w-full px-4 py-3.5 border-2 rounded-xl text-sm transition-all duration-300 outline-none font-mono tracking-widest pr-12"
                    :class="{
                      'border-green-400 bg-green-50/60': siretStatus === 'found',
                      'border-red-300 bg-red-50/60': siretStatus === 'error',
                      'border-[#e2e8f0] focus:border-[#006d35]': !siretStatus
                    }"
                  />
                  <div class="absolute right-4 top-1/2 -translate-y-1/2 w-6 h-6 flex items-center justify-center">
                    <div v-if="siretLoading" class="w-5 h-5 border-2 border-[#006d35] border-t-transparent rounded-full animate-spin"></div>
                    <CheckCircleIcon v-else-if="siretStatus === 'found'" class="w-5 h-5 text-green-500" />
                    <XCircleIcon v-else-if="siretStatus === 'error'" class="w-5 h-5 text-red-400" />
                    <span v-else class="text-xs text-gray-300 font-mono tabular-nums">{{ siretInput.length }}/14</span>
                  </div>
                </div>
              </div>

              <Transition name="card-pop">
                <div v-if="companyData" class="mb-4 rounded-2xl border-2 border-green-200 overflow-hidden" style="background:linear-gradient(135deg,#f0fdf4,#dcfce7);">
                  <div class="p-4">
                    <div class="flex items-start gap-3">
                      <div class="w-11 h-11 rounded-xl bg-white shadow-sm flex items-center justify-center shrink-0">
                        <BuildingOfficeIcon class="w-6 h-6 text-gray-400" />
                      </div>
                      <div class="flex-1 min-w-0">
                        <p class="font-bold text-[#001d32] text-base leading-tight">
                          {{ companyData.infolegales.nomcommercialrne || companyData.infolegales.nomcommercialinsee || companyData.infolegales.denorne || companyData.infolegales.denoinsee }}
                        </p>
                        <p class="text-xs text-gray-500 mt-0.5 truncate">
                          {{ [companyData.infolegales.voieadressageinsee, companyData.infolegales.codepostalinsee, companyData.infolegales.villeinsee].filter(Boolean).join(', ') || t('public.register.siretAddressMissing') }}
                        </p>
                        <div class="flex flex-wrap gap-1.5 mt-2">
                          <span v-if="companyData.infolegales.naflibinsee || companyData.infolegales.nafinsee" class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-green-200 text-green-700 font-medium">
                            {{ companyData.infolegales.naflibinsee || companyData.infolegales.nafinsee }}
                          </span>
                          <span v-if="companyData.infolegales.catjurlibinsee" class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-gray-200 text-gray-500">
                            {{ companyData.infolegales.catjurlibinsee }}
                          </span>
                          <span class="text-xs px-2 py-0.5 rounded-full bg-white/80 border border-gray-200 text-gray-500 font-mono">
                            {{ t('public.register.siretSiretBadge', { siret: siretInput }) }}
                          </span>
                        </div>
                      </div>
                    </div>
                  </div>
                  <button
                    @click="useCompanyData"
                    class="w-full py-3 text-white text-sm font-semibold transition hover:opacity-90 flex items-center justify-center gap-2"
                    style="background:#006d35;"
                  >
                    <SparklesIcon class="w-4 h-4" />
                    {{ t('public.register.siretUseInfo') }}
                  </button>
                </div>
              </Transition>

              <Transition name="fade">
                <div v-if="siretStatus === 'error'" class="mb-4 p-3 rounded-xl bg-red-50 border border-red-200 flex items-start gap-2 text-sm text-red-600">
                  <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" />
                  <span>{{ t('public.register.siretError') }}</span>
                </div>
              </Transition>

              <div class="pt-4 border-t border-[#e2e8f0] text-center">
                <button @click="skipSiret" class="text-sm text-[#40617f] hover:text-[#001d32] transition hover:underline underline-offset-2">
                  {{ t('public.register.siretSkip') }}
                </button>
              </div>
            </div>
          </div>

          <div v-else-if="step === 'siret-confirm'" key="siret-confirm" class="max-w-[640px] mx-auto w-full">
            <button @click="step = 'siret'" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#001d32] transition mb-6">
              <ArrowLeftIcon class="w-4 h-4" /> {{ t('public.register.back') }}
            </button>

            <div class="flex flex-col items-center text-center mb-8">
              <p class="text-[#006d35] text-xs font-bold uppercase tracking-widest mb-3">{{ t('public.register.siretConfirmStep') }}</p>
              <h2 class="font-jakarta font-extrabold text-[#001d32] text-4xl tracking-tight">{{ t('public.register.siretConfirmTitle') }}</h2>
              <p class="text-[#40617f] text-base mt-3 max-w-md">{{ t('public.register.siretConfirmSubtitle') }}</p>
            </div>

            <div v-if="companyData" class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] overflow-hidden">

              <div class="p-8 flex items-start gap-6 border-b border-[#edf4ff]">
                <div class="w-20 h-20 rounded-[20px] bg-[#edf4ff] flex items-center justify-center shrink-0">
                  <BuildingOfficeIcon class="w-10 h-10 text-[#40617f]" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="flex items-center gap-3 flex-wrap">
                    <h3 class="font-jakarta font-bold text-[#001d32] text-2xl">
                      {{ companyData.infolegales.nomcommercialrne || companyData.infolegales.nomcommercialinsee || companyData.infolegales.denorne || companyData.infolegales.denoinsee }}
                    </h3>
                    <span class="bg-[#d1fae5] text-[#006d35] font-bold text-xs uppercase tracking-widest px-3 py-1 rounded-full flex items-center gap-1">
                      <CheckBadgeIcon class="w-3.5 h-3.5" /> {{ t('public.register.siretActiveBadge') }}
                    </span>
                  </div>
                  <p class="text-[#40617f] text-sm mt-1">
                    {{ [companyData.infolegales.voieadressageinsee, companyData.infolegales.codepostalinsee, companyData.infolegales.villeinsee].filter(Boolean).join(', ') || t('public.register.siretAddressMissing') }}
                  </p>
                </div>
              </div>

              <div class="p-8 grid grid-cols-2 gap-4">
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.register.siretFieldSiret') }}</p>
                  <p class="font-mono font-bold text-[#001d32] text-sm">{{ siretInput }}</p>
                </div>
                <div class="bg-[#f7f9ff] rounded-xl p-4">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.register.siretFieldNafCode') }}</p>
                  <p class="font-bold text-[#001d32] text-sm">{{ companyData.infolegales.nafinsee || '—' }}</p>
                </div>
                <div class="bg-[#f7f9ff] rounded-xl p-4 col-span-2">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.register.siretFieldActivity') }}</p>
                  <p class="font-bold text-[#001d32] text-sm">{{ companyData.infolegales.naflibinsee || companyData.infolegales.nafinsee || '—' }}</p>
                </div>
                <div v-if="companyData.infolegales.catjurlibinsee" class="bg-[#f7f9ff] rounded-xl p-4 col-span-2">
                  <p class="text-[#40617f] text-xs uppercase tracking-wide mb-1">{{ t('public.register.siretFieldLegalForm') }}</p>
                  <p class="font-bold text-[#001d32] text-sm">{{ companyData.infolegales.catjurlibinsee }}</p>
                </div>
              </div>

              <div class="px-8 pb-8 flex flex-col gap-3">
                <button
                  @click="useCompanyData"
                  class="w-full py-4 rounded-xl text-white font-jakarta font-bold text-base transition hover:opacity-90 flex items-center justify-center gap-2"
                  style="background: linear-gradient(135deg, #006d35 0%, #1b8848 100%);"
                >
                  <CheckIcon class="w-5 h-5" />
                  {{ t('public.register.siretConfirmBtn') }}
                </button>
                <button
                  @click="step = 'siret'"
                  class="w-full py-3 rounded-xl text-[#40617f] font-medium text-sm hover:text-[#001d32] transition"
                >
                  {{ t('public.register.siretRefuseBtn') }}
                </button>
              </div>
            </div>
          </div>

          <div v-else-if="step === 'form'" key="form" class="max-w-lg mx-auto w-full">
            <button @click="goBack" class="inline-flex items-center gap-1 text-sm text-[#40617f] hover:text-[#001d32] transition mb-6">
              <ArrowLeftIcon class="w-4 h-4" /> {{ t('public.register.back') }}
            </button>

            <div class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] p-10">
              <div class="flex items-center gap-4 mb-6">
                <div class="w-12 h-12 rounded-xl flex items-center justify-center"
                  :style="accountType === 'provider' ? 'background:#edf4ff;' : 'background:#d8eaff;'">
                  <WrenchScrewdriverIcon v-if="accountType === 'provider'" class="w-5 h-5 text-[#40617f]" />
                  <UserIcon v-else class="w-5 h-5 text-[#40617f]" />
                </div>
                <div>
                  <h2 class="font-jakarta font-bold text-[#001d32] text-xl">
                    {{ accountType === 'provider' ? t('public.register.formProviderTitle') : t('public.register.formIndividualTitle') }}
                  </h2>
                  <p v-if="companyData" class="text-sm font-medium text-[#006d35] flex items-center gap-1">
                    <SparklesIcon class="w-3.5 h-3.5" /> {{ t('public.register.formPrefilledNotice') }}
                  </p>
                  <p v-else class="text-sm text-[#40617f]">{{ t('public.register.formFillFields') }}</p>
                </div>
              </div>

              <form @submit.prevent="handleRegister" class="space-y-4">

                <template v-if="accountType === 'provider'">
                  <div class="rounded-2xl border border-[#e2e8f0] p-4 space-y-3 bg-[#f7f9ff]">
                    <p class="text-xs font-bold text-[#40617f] uppercase tracking-wider flex items-center gap-1.5">
                      <BuildingOfficeIcon class="w-3.5 h-3.5" /> {{ t('public.register.sectionCompany') }}
                    </p>
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldCompanyName') }} <span class="text-red-400">*</span></label>
                      <input v-model="form.companyName" type="text" required
                        class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors"
                        :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                      />
                    </div>
                    <div class="grid grid-cols-2 gap-3">
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldSiret') }}</label>
                        <input v-model="form.siret" type="text" maxlength="14"
                          class="w-full px-3 py-2.5 border rounded-xl text-xs focus:outline-none transition-colors font-mono"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldSiretPlaceholder')"
                        />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldActivity') }}</label>
                        <input v-model="form.activity" type="text"
                          class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldActivityPlaceholder')"
                        />
                      </div>
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldAddress') }}</label>
                      <AddressAutocomplete
                        :model-value="[form.streetNumber, form.streetName, form.postalCode, form.city].filter(Boolean).join(' ')"
                        @update:model-value="() => {}"
                        @select="handleAddressSelect"
                        :placeholder="t('public.register.fieldAddressSearch')"
                        class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors pr-10"
                        :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                      />
                    </div>
                    <div class="grid grid-cols-[80px_1fr] gap-2">
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldStreetNumber') }}</label>
                        <input v-model="form.streetNumber" type="text"
                          class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldStreetNumberPlaceholder')" />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldStreetName') }}</label>
                        <input v-model="form.streetName" type="text"
                          class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldStreetNamePlaceholder')" />
                      </div>
                    </div>
                    <div class="grid grid-cols-[120px_1fr] gap-2">
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldPostalCode') }}</label>
                        <input v-model="form.postalCode" type="text" maxlength="5"
                          class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors font-mono"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldPostalCodePlaceholder')" />
                      </div>
                      <div>
                        <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldCity') }}</label>
                        <input v-model="form.city" type="text"
                          class="w-full px-3 py-2.5 border rounded-xl text-sm focus:outline-none transition-colors"
                          :class="companyData ? 'border-green-300 bg-green-50/40' : 'border-[#e2e8f0] bg-white focus:border-[#006d35]'"
                          :placeholder="t('public.register.fieldCityPlaceholder')" />
                      </div>
                    </div>
                  </div>
                </template>

                <div class="space-y-3">
                  <p v-if="accountType === 'provider'" class="text-xs font-bold text-[#40617f] uppercase tracking-wider pt-1 flex items-center gap-1.5">
                    <UserIcon class="w-3.5 h-3.5" /> {{ t('public.register.sectionContact') }}
                  </p>
                  <div class="grid grid-cols-2 gap-3">
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldFirstName') }} <span class="text-red-400">*</span></label>
                      <input v-model="form.firstName" type="text" required
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white" />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldLastName') }} <span class="text-red-400">*</span></label>
                      <input v-model="form.lastName" type="text" required
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white" />
                    </div>
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldEmail') }} <span class="text-red-400">*</span></label>
                    <input v-model="form.email" type="email" required
                      class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white" />
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldPhone') }}</label>
                    <input v-model="form.phone" type="tel"
                      class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white"
                      :placeholder="t('public.register.fieldPhonePlaceholder')" />
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldAddress') }}</label>
                    <AddressAutocomplete
                      :model-value="[form.streetNumber, form.streetName, form.postalCode, form.city].filter(Boolean).join(' ')"
                      @update:model-value="() => {}"
                      @select="handleAddressSelect"
                      :placeholder="t('public.register.fieldAddressSearch')"
                      class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white pr-10"
                    />
                  </div>
                  <div class="grid grid-cols-[80px_1fr] gap-2">
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldStreetNumber') }}</label>
                      <input v-model="form.streetNumber" type="text"
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white"
                        :placeholder="t('public.register.fieldStreetNumberPlaceholder')" />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldStreetName') }}</label>
                      <input v-model="form.streetName" type="text"
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white"
                        :placeholder="t('public.register.fieldStreetNamePlaceholder')" />
                    </div>
                  </div>
                  <div class="grid grid-cols-[120px_1fr] gap-2">
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldPostalCode') }}</label>
                      <input v-model="form.postalCode" type="text" maxlength="5"
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white font-mono"
                        :placeholder="t('public.register.fieldPostalCodePlaceholder')" />
                    </div>
                    <div>
                      <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldCity') }}</label>
                      <input v-model="form.city" type="text"
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white"
                        :placeholder="t('public.register.fieldCityPlaceholder')" />
                    </div>
                  </div>
                  <div>
                    <label class="block text-xs font-medium text-[#001d32] mb-1">{{ t('public.register.fieldPassword') }} <span class="text-red-400">*</span></label>
                    <div class="relative">
                      <input v-model="form.password" :type="showPassword ? 'text' : 'password'" required minlength="8"
                        class="w-full px-3 py-2.5 border border-[#e2e8f0] rounded-xl text-sm focus:outline-none focus:border-[#006d35] transition-colors bg-white pr-10" />
                      <button type="button" @click="showPassword = !showPassword" class="absolute right-3 top-1/2 -translate-y-1/2 text-gray-400 hover:text-gray-600 transition">
                        <EyeSlashIcon v-if="showPassword" class="w-5 h-5" />
                        <EyeIcon v-else class="w-5 h-5" />
                      </button>
                    </div>
                    <div v-if="form.password" class="mt-2 flex gap-1">
                      <div v-for="i in 4" :key="i" class="h-1.5 flex-1 rounded-full transition-all duration-300"
                        :style="passwordStrength >= i ? `background:${strengthColor}` : 'background:#e2e8f0'"></div>
                    </div>
                    <p v-if="form.password" class="text-xs mt-1 transition-all" :style="`color:${strengthColor}`">{{ strengthLabel }}</p>
                  </div>
                </div>

                <Transition name="fade">
                  <p v-if="error" class="text-red-600 text-sm bg-red-50 p-3 rounded-xl border border-red-200 flex items-start gap-2">
                    <ExclamationTriangleIcon class="w-4 h-4 shrink-0 mt-0.5" /><span>{{ error }}</span>
                  </p>
                </Transition>

                <button type="submit" :disabled="loading"
                  class="w-full py-4 rounded-xl font-jakarta font-bold text-white text-base transition-all duration-200 hover:opacity-90 active:scale-[0.99] disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center gap-2"
                  style="background: linear-gradient(135deg, #006d35, #1b8848);"
                >
                  <div v-if="loading" class="w-4 h-4 border-2 border-white border-t-transparent rounded-full animate-spin"></div>
                  <CheckCircleIcon v-else class="w-4 h-4" />
                  <span>{{ loading ? t('public.register.submitting') : t('public.register.submit') }}</span>
                </button>
              </form>
            </div>

            <p class="text-center text-sm text-[#40617f] mt-6">
              {{ t('public.register.alreadyRegistered') }}
              <RouterLink to="/connexion" class="font-semibold hover:underline underline-offset-2 text-[#006d35]">{{ t('public.register.loginLink') }}</RouterLink>
            </p>
          </div>

          <div v-else-if="step === 'done'" key="done" class="max-w-lg mx-auto w-full text-center">
            <div class="bg-white rounded-[24px] shadow-[0px_12px_40px_0px_rgba(0,29,50,0.06)] p-10 space-y-6">
              <div class="w-20 h-20 rounded-full flex items-center justify-center mx-auto" style="background:#dcfce7;">
                <EnvelopeIcon class="w-10 h-10 text-[#006d35]" />
              </div>
              <div>
                <h2 class="font-jakarta font-extrabold text-[#001d32] text-2xl mb-2">{{ t('public.register.doneTitle') }}</h2>
                <p class="text-[#40617f] text-base leading-relaxed">
                  {{ t('public.register.doneText') }} <strong class="text-[#001d32]">{{ form.email }}</strong>.
                </p>
                <p class="text-[#40617f] text-sm mt-2">
                  {{ t('public.register.doneHint') }}
                </p>
              </div>
              <p class="text-xs text-gray-400">{{ t('public.register.doneSpam') }}</p>
              <RouterLink to="/connexion" class="inline-flex items-center gap-2 text-sm font-semibold text-[#40617f] hover:text-[#006d35] transition">
                {{ t('public.register.doneBack') }}
              </RouterLink>
            </div>
          </div>

        </Transition>
      </div>
    </main>

    <footer class="border-t border-[#edf4ff]">
      <div class="max-w-[1280px] mx-auto px-8 py-12 flex flex-col sm:flex-row items-center justify-between gap-4">
        <div>
          <p class="font-jakarta font-bold text-[#006d35] text-xl">UpcycleConnect</p>
          <p class="text-[#40617f] text-sm mt-1">{{ t('public.register.footerCopyright') }}</p>
        </div>
        <div class="flex gap-8 text-sm text-[#40617f]">
          <a href="#" class="hover:text-[#006d35] transition">{{ t('public.register.footerPrivacy') }}</a>
          <a href="#" class="hover:text-[#006d35] transition">{{ t('public.register.footerTerms') }}</a>
          <a href="#" class="hover:text-[#006d35] transition">{{ t('public.register.footerCommitment') }}</a>
        </div>
      </div>
    </footer>

  </div>
</template>

<script setup>import { ref, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useUserAuthStore } from '@/stores/userAuth'
import AddressAutocomplete from '@/components/AddressAutocomplete.vue'
import {
  UserIcon,
  WrenchScrewdriverIcon,
  CheckIcon,
  CheckCircleIcon,
  CheckBadgeIcon,
  XCircleIcon,
  MagnifyingGlassIcon,
  BuildingOfficeIcon,
  EyeIcon,
  EyeSlashIcon,
  SparklesIcon,
  ExclamationTriangleIcon,
  ArrowLeftIcon,
  ChevronRightIcon,
  LockClosedIcon,
  EnvelopeIcon,
} from '@heroicons/vue/24/outline'

const { t } = useI18n()
const router   = useRouter()
const userAuth = useUserAuthStore()

const step = ref('type')
const accountType = ref('')
const siretInput = ref('')
const siretLoading = ref(false)
const siretStatus = ref('')
const companyData = ref(null)
const showPassword = ref(false)
const loading = ref(false)
const error = ref('')

const form = ref({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  password: '',
  companyName: '',
  siret: '',
  activity: '',
  streetNumber: '',
  streetName: '',
  postalCode: '',
  city: '',
})

const stepIndex = computed(() => ({ type: 0, siret: 1, 'siret-confirm': 1, form: accountType.value === 'provider' ? 2 : 1 }[step.value] ?? 0))

function goToNext() {
  if (accountType.value === 'provider') {
    step.value = 'siret'
  } else {
    step.value = 'form'
  }
}

function goBack() {
  step.value = accountType.value === 'provider' ? 'siret' : 'type'
}

function skipSiret() {
  companyData.value = null
  step.value = 'form'
}

let siretTimer = null
function onSiretInput() {
  const digits = siretInput.value.replace(/\D/g, '').slice(0, 14)
  siretInput.value = digits
  siretStatus.value = ''
  companyData.value = null
  clearTimeout(siretTimer)
  if (digits.length === 14) {
    siretTimer = setTimeout(() => fetchSiret(digits), 400)
  }
}

async function fetchSiret(siret) {
  siretLoading.value = true
  siretStatus.value = ''
  try {
    const res = await fetch(`/api/public/v1/siret/${siret}`)
    if (res.status === 404) { siretStatus.value = 'error'; return }
    if (!res.ok) throw new Error('server_error')
    const data = await res.json()
    if (data?.infolegales) {
      companyData.value = data
      siretStatus.value = 'found'
    } else {
      siretStatus.value = 'error'
    }
  } catch {
    siretStatus.value = 'error'
  } finally {
    siretLoading.value = false
  }
}

function handleAddressSelect(props) {
  form.value.streetNumber = props.housenumber || ''
  form.value.streetName   = props.street || props.name || ''
  form.value.postalCode   = props.postcode || ''
  form.value.city         = props.city || ''
}

function useCompanyData() {
  const d = companyData.value?.infolegales
  if (!d) return
  if (step.value === 'siret-confirm') {
    form.value.companyName  = d.nomcommercialrne || d.nomcommercialinsee || d.denorne || d.denoinsee || ''
    form.value.siret        = siretInput.value
    form.value.activity     = d.naflibinsee || d.nafinsee || ''
    const rawAddress        = (d.voieadressageinsee || '').trim()
    const numMatch          = rawAddress.match(/^(\d+[A-Za-z]?)\s+(.+)$/)
    form.value.streetNumber = numMatch ? numMatch[1] : ''
    form.value.streetName   = numMatch ? numMatch[2] : rawAddress
    form.value.postalCode   = d.codepostalinsee || ''
    form.value.city         = d.villeinsee || ''
    step.value = 'form'
  } else {
    step.value = 'siret-confirm'
  }
}

const passwordStrength = computed(() => {
  const p = form.value.password
  if (!p) return 0
  let s = 0
  if (p.length >= 8)          s++
  if (/[A-Z]/.test(p))        s++
  if (/[0-9]/.test(p))        s++
  if (/[^A-Za-z0-9]/.test(p)) s++
  return s
})

const strengthColor = computed(() => ['#e2e8f0','#ef4444','#f97316','#eab308','#22c55e'][passwordStrength.value])
const strengthLabel = computed(() => [
  '',
  t('public.register.strengthVeryWeak'),
  t('public.register.strengthWeak'),
  t('public.register.strengthMedium'),
  t('public.register.strengthStrong'),
][passwordStrength.value])

async function handleRegister() {
  loading.value = true
  error.value = ''
  try {
    const addressParts = [
      form.value.streetNumber,
      form.value.streetName,
      form.value.postalCode,
      form.value.city,
    ].map(s => s.trim()).filter(Boolean)
    const fullAddress = addressParts.length ? addressParts.join(' ') : null

    await userAuth.register({
      first_name:   form.value.firstName,
      last_name:    form.value.lastName,
      email:        form.value.email,
      password:     form.value.password,
      phone:        form.value.phone || null,
      account_type: accountType.value,
      company_name: form.value.companyName || null,
      siret:        form.value.siret || null,
      activity:     form.value.activity || null,
      address:      fullAddress,
    })
    step.value = 'done'
  } catch (e) {
    error.value = e.message || t('public.register.errorGeneric')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.step-enter-active,
.step-leave-active {
  transition: opacity 0.2s ease, transform 0.25s ease;
}
.step-enter-from { opacity: 0; transform: translateX(20px); }
.step-leave-to   { opacity: 0; transform: translateX(-20px); }

.card-pop-enter-active { transition: all 0.4s cubic-bezier(0.34, 1.56, 0.64, 1); }
.card-pop-leave-active { transition: all 0.2s ease; }
.card-pop-enter-from, .card-pop-leave-to { opacity: 0; transform: scale(0.95) translateY(-4px); }

.badge-enter-active { transition: all 0.3s cubic-bezier(0.34, 1.56, 0.64, 1); }
.badge-leave-active { transition: all 0.15s ease; }
.badge-enter-from, .badge-leave-to { opacity: 0; transform: scale(0.5); }

.fade-enter-active, .fade-leave-active { transition: opacity 0.2s ease; }
.fade-enter-from, .fade-leave-to       { opacity: 0; }
</style>
