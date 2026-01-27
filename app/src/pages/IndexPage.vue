<template>

    <q-card class="w-full max-w-5xl rounded-xl shadow-lg">

      <!-- Header -->
      <q-card-section class="bg-primary text-white rounded-t-xl">
        <div class="text-lg font-semibold">
          Teste de API
        </div>
      </q-card-section>

      <!-- Form -->
      <q-card-section class="space-y-4">
        <div class="grid grid-cols-12 gap-4">

          <!-- Método -->
          <q-select
            v-model="method"
            :options="methodOptions"
            label="Método"
            filled
            dense
            class="col-span-12 md:col-span-2"
          />

          <!-- URL -->
          <q-input
            v-model="urlApi"
            label="Rota da API"
            filled
            dense
            placeholder="/v1/produtos"
            class="col-span-12 md:col-span-10"
          />

          <!-- Payload -->
          <q-input
            v-model="payload"
            type="textarea"
            autogrow
            filled
            label="Payload (JSON)"
            placeholder="{ }"
            class="col-span-12 font-mono"
          />
        </div>

        <!-- Botão -->
        <div class="flex justify-end">
          <q-btn
            color="primary"
            icon="send"
            label="Executar"
            no-caps
            class="px-6"
            @click="fetchApi"
          />
        </div>
      </q-card-section>

      <!-- Response -->
      <q-separator />

      <q-card-section>
        <div class="text-sm font-semibold text-gray-600 mb-2">
          Resposta
        </div>

        <q-card
          flat
          bordered
          class=" rounded-lg"
          :class="{
            'text-green-400': !isError,
            'text-red-400': isError

          }"
        >
          <q-card-section class="p-3 overflow-auto max-h-96">
            <pre class="text-xs leading-relaxed"> {{ JSON.stringify(data, null, 2) }} </pre>
          </q-card-section>
        </q-card>
      </q-card-section>

    </q-card>
</template>


<script setup lang="ts">
    import { api } from 'src/boot/axios';
    import { ref } from 'vue';

    const methodOptions: string[] = [
        "GET",
        "POST",
        "PUT",
        "DELETE"
    ];

    const data = ref<string>('');
    const urlApi = ref<string>('');
    const method = ref<'GET'|'POST'|'PUT'|'DELETE'>('GET');
    const payload = ref<string>();
    let isError = ref<boolean>(false);

    const fetchApi = async () => {
        console.log(method.value);

        data.value = '';
        let res;
        try {
            switch (method.value) {
                case 'GET':
                    res = await api.get(`/${urlApi.value}`)
                    data.value = res.data;

                    break;

                case 'POST':
                    res = await api.post(`/${urlApi.value}`, payload.value);
                    data.value = res.data;

                    break;
                case 'PUT':
                    res = await api.put(`/${urlApi.value}`, payload.value);
                    data.value = res.data;

                    break;

                case 'DELETE':
                    res = await api.delete(`/${urlApi.value}`);
                    data.value = res.data;

                    break;

                default:
                    break;
            }
            isError.value = false;

        } catch (error) {
            console.error('Erro: ', error);

            data.value = error.response?.data;
            isError.value = true;
        }
    };

</script>
