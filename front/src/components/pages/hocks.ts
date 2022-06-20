import { useState } from "react";

export const useWallet = () => {
  const [publicKey, setPublicKey] = useState<string>("");
  const [privateKey, setPrivateKey] = useState<string>("");
  const [blockchainAddress, sebBlockchainAddress] = useState<string>("");

  const createWallet = async () => {
    const res = await fetch("/api/blockchain/wallet", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
    });
    if (res.ok) {
      const wallet = await res.json();
      setPublicKey(wallet.publicKey);
      setPrivateKey(wallet.privateKey);
      sebBlockchainAddress(wallet.blockchainAddress);
    }
  };

  return {
    publicKey,
    privateKey,
    blockchainAddress,
    createWallet,
  };
};

export const useAmount = () => {
  const [currentAmount, setCurrentAmount] = useState<string>("");

  const getAmount = async (blockchainAddress: string) => {
    try {
      const query = new URLSearchParams({
        blockchainAddress,
      });
      const res = await fetch(`/api/blockchain/wallet?${query}`, {
        method: "GET",
        headers: {
          "Content-Type": "application/json",
        },
      });
      if (res.ok) {
        const wallet = await res.json();
        setCurrentAmount(wallet.amount);
      }
    } catch (e) {
      console.error(e);
    }
  };

  return {
    currentAmount,
    getAmount,
  };
};

export const useSendMoney = () => {
  const [recipientKey, setRecipientKey] = useState<string>("");
  const [amount, setAmount] = useState<string>("");

  const sendMoney = async (
    senderPrivateKey: string,
    senderPublicKey: string,
    senderBlockchainAddress: string
  ) => {
    await fetch("/api/blockchain/transactions", {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        senderPrivateKey,
        senderBlockchainAddress,
        recipientBlockchainAddress: recipientKey,
        senderPublicKey,
        value: amount,
      }),
    });
  };
  return {
    recipientKey,
    setRecipientKey,
    amount,
    setAmount,
    sendMoney,
  };
};
