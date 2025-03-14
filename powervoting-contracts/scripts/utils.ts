const fs = require("fs");
const path = require("path");
const fse = require("fs-extra");
import hre from "hardhat";

const CONSTANT_PATH = path.join(__dirname, hre.network.name + "_config.json");

const getConstantJson = function () {
  let data;
  if (fs.existsSync(CONSTANT_PATH)) {
    data = fse.readJsonSync(CONSTANT_PATH);
  } else {
    data = {};
  }
  return data;
};
const getConstantValue = async function (key: string) {
  let data;
  if (fs.existsSync(CONSTANT_PATH)) {
    data = fse.readJsonSync(CONSTANT_PATH);
    console.log(data);
  } else {
    data = {};
  }
  return data[key];
};

const updateConstant = async function (key: string, value: string) {
  let data;
  if (fs.existsSync(CONSTANT_PATH)) {
    data = fse.readJsonSync(CONSTANT_PATH);
  } else {
    data = {};
  }
  data[key] = value;

  fse.writeJsonSync(CONSTANT_PATH, data);
};

export { getConstantValue, updateConstant, getConstantJson, CONSTANT_PATH };
