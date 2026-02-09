import { useQuasar } from "quasar";

export type NotifyType = 'positive' | 'negative' | 'warning' | 'info';

export function useNotify() {
    const $q = useQuasar();

    const notify = (type: NotifyType, message: string) => {
        $q.notify({
            type: type,
            position: 'top',
            message: message,
        });
    };

    return { notify };
};
