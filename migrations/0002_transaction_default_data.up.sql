INSERT INTO `walletapp`.`users` (`name`, `balance`) VALUES ('bismo', '100000');
INSERT INTO `walletapp`.`users` (`name`, `balance`) VALUES ('andre', '20000');
INSERT INTO `walletapp`.`users` (`name`, `balance`) VALUES ('test', '50000');

INSERT INTO `walletapp`.`transactions` (`transaction_type`, `amount`, `user_id`) VALUES ('DEBIT', '50000', 1);
INSERT INTO `walletapp`.`transactions` (`transaction_type`, `amount`, `user_id`) VALUES ('DEBIT', '20000', 2);
INSERT INTO `walletapp`.`transactions` (`transaction_type`, `amount`, `user_id`) VALUES ('DEBIT', '50000', 3);
INSERT INTO `walletapp`.`transactions` (`transaction_type`, `amount`, `user_id`) VALUES ('CREDIT', '30000', 1);
INSERT INTO `walletapp`.`transactions` (`transaction_type`, `amount`, `user_id`) VALUES ('DEBIT', '80000', 1);
