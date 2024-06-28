import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
// Chinese language pack
import zh from './zh.json';
// English language pack
import en from './en.json';
 
const resources = {
  en: {
    translation: en
  },
  zh: {
    translation: zh
  }
};
 
i18n.use(initReactI18next).init({
  resources,
  lng: 'en', //Set default language
  interpolation: {
    escapeValue: false
  }
});
 
export default i18n;