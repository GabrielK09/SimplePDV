<template>
    <q-dialog v-model="internalDialog" persistent>
        <div class="flex flex-col bg-white w-full phone:bg-black shadow-lg">

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

    const generateQrCodePix = async () => {
        const base64 = await generateQRCodeBuilder(props.totalSale, props.pixKey);
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
        const errorSaleNotify = notify(
            'negative',
            'Erro ao processar a venda, tente novamente.'
        );
        if(props.totalSale < 0)
        {
            errorSaleNotify;
        } else if(props.pixKey === ''){
            errorSaleNotify;
        } else {
            generateQrCodePix();
        };

        console.log('Oie');

    });
</script>
