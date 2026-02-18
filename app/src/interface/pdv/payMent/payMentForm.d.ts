interface PayMentFormContract {
    readonly id: number;
    specie: string;
    pix_key: string;
}; // Para create
interface PayMentPayLoadContract {
    sale_id: number;
    specie: string;
    amount_paid: number;
};

interface PayMentValue {
    id: number;
    specie: string;
    amount: string;
    pix_key: string;
}; // Para venda
