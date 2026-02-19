/**
 * Para create
 */
interface PayMentFormContract {
    readonly id: number;
    specie: string;
    pix_key: string;
};

/**
 * Para o pagamento na api
 */
interface PayMentPayLoadContract {
    sale_id: number;
    species: any[];

};

/**
 * Para venda
 */
interface PayMentValue {
    id: number;
    specie: string;
    amount: string;
    pix_key: string;
};
