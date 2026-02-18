/**
 * Para create
 */
interface PayMentFormContract {
    readonly id: number;
    specie: string;
    pix_key: string;
};

type Species = {
    specie: string;
    amount_paid: number;
};

/**
 * Para o pagamento
 */
interface PayMentPayLoadContract {
    sale_id: number;
    species: Species;

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
