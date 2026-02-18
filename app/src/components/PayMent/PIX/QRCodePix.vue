<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="flex flex-col bg-white phone:bg-black shadow-lg px-4">
            <h2 class="text-gray-600 text-center">
                Pagamento via PIX
            </h2>

            <span class="text-gray-600 text-center">
                R$ {{ props.totalSale.toFixed(2).toString().replace('.', ',') }}
            </span>

            <div>
                <img :src=imgBase64Pix />

            </div>

            <div class="my-4 flex flex-center">
                <q-btn
                    color="red"
                    flat
                    label="Cancelar"
                    no-caps
                    class="mr-2"
                    @click="emits('close', true)"
                />

                <q-btn
                    color="primary"
                    label="Finalizar"
                    no-caps
                    class="ml-2"
                    @click="emits('confirm', true)"
                />
            </div>
        </div>
    </q-dialog>
</template>

<script setup lang="ts">
    import { useNotify } from 'src/helpers/QNotify/useNotify';
    import { generateQRCodeBuilder } from 'src/services/QRCode/QrCodePix';
    import { onMounted, ref } from 'vue';

    const internalDialog = ref<boolean>(true);
    const imgBase64Pix = ref<string>('');
    const { notify } = useNotify();

    const props = defineProps<{
        totalSale: number,
        pixKey: string
    }>();

    const emits = defineEmits<{
        (e: 'close', value: boolean),
        (e: 'confirm', value: boolean)
    }>();

    const generateQrCodePix = async () => {
        const res = await generateQRCodeBuilder(props.totalSale, props.pixKey);
        const base64 = res.base64;

        console.log(res);

        if(!base64)
        {
            notify(
                'negative',
                'Erro ao gerar o QR code, tente novamente.'
            );

            internalDialog.value = false;
            return;

        };

        imgBase64Pix.value = base64;
    };

    onMounted(() => {
        if(props.totalSale < 0)
        {
            notify(
                'negative',
                'Erro ao processar a venda, tente novamente.'
            );

        } else if(props.pixKey === ''){
            notify(
                'negative',
                'Erro ao processar a venda, tente novamente.'
            );

        } else {
            generateQrCodePix();

        };
    });
</script>
