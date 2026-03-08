enum ReportTypes {
    CashRegister = "cash-register",
    PayMentForms = "pay-ment-forms",
    SaledItens = "saled-itens"
};

interface ReportContract {
    report_type: ReportTypes|null;
    start_date: string|null;
    end_date: string|null;
};
