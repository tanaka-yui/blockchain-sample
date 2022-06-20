import Button from "@mui/material/Button";
import TextField from "@mui/material/TextField";
import React, { useEffect } from "react";
import { useAmount, useSendMoney, useWallet } from "~/components/pages/hocks";
import { Box, Container } from "@mui/material";

const Index: React.FC = () => {
  const { publicKey, privateKey, blockchainAddress, createWallet } =
    useWallet();

  const { recipientKey, setRecipientKey, amount, setAmount, sendMoney } =
    useSendMoney();

  const { currentAmount, getAmount } = useAmount();

  const callCreateWallet = async () => {
    try {
      const res = await createWallet();
      console.log(res);
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    callCreateWallet();
    return () => {};
  }, []);

  useEffect(() => {
    const intervalId = setInterval(async () => {
      await getAmount(blockchainAddress);
    }, 5000);
    getAmount(blockchainAddress);
    return () => {
      clearTimeout(intervalId);
    };
  }, [blockchainAddress]);

  return (
    <Container>
      <h1>Wallet</h1>
      <h2>Public Key</h2>
      <p>{publicKey}</p>
      <h2>Private Key</h2>
      <p>{privateKey}</p>
      <h2>BlockChain Address</h2>
      <p>{blockchainAddress}</p>
      <h2>Wallet Amount</h2>
      <p>{currentAmount || 0}</p>
      <hr />
      <Box>
        <h2>Send Money</h2>
        <h3>Recipient Key</h3>
        <TextField
          value={recipientKey}
          fullWidth
          onChange={(e) => setRecipientKey(e.target.value)}
        />
        <h3>Amount</h3>
        <TextField
          value={amount}
          fullWidth
          onChange={(e) => setAmount(e.target.value)}
        />
      </Box>
      <hr />
      <Button
        variant="outlined"
        onClick={() => sendMoney(privateKey, publicKey, blockchainAddress)}
      >
        Send
      </Button>
    </Container>
  );
};

export default Index;
