const Maze = artifacts.require("Maze");

module.exports = async function (deployer) {
  await deployer.deploy(Maze);
};
