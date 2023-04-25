const spinProccessor = artifacts.require("SpinProccessor");

module.exports = function (_deployer, network, accounts) {
  console.log("accounts", accounts);
  const owner = accounts[0];
  const numberConbfirmationRequired = 2;
  const usdtAddr = "0xD0355200111C2B21AAbC1a31552eCCDc5d4E905d";
  _deployer.deploy(spinProccessor, usdtAddr);
};
