type LocalStorageObj = {
  [key: string]: string,
};

async function getLocalStorage(key: string): Promise<string> {
  return chrome.storage.local.get(key)
    .then((result) => {
      if (result === undefined || result[key] === undefined) {
        return '';
      }
      return result[key] as string;
    })
    .catch((e :unknown) => {
      // TODO: Some generic error catching here. Could be improved.
      console.error(e);
      return '';
    });
}

async function setLocalStorage(key: string, value:string): Promise<void> {
  const storedObj: LocalStorageObj = {};
  storedObj[key] = value;
  chrome.storage.local.set(storedObj);
}

async function localStorageExists(key: string): Promise<boolean> {
  return chrome.storage.local.get(key)
    .then((result) => {
      if (result === undefined || result[key] !== undefined) {
        return false;
      }
      const resValue = result[key] as string;
      return resValue !== '';
    })
    .catch((e: unknown) => {
      // TODO: Some generic error catching here. Could be improved.
      console.error(e);
      return false;
    });
}

async function deleteLocalStorage(key: string): Promise<boolean> {
  return chrome.storage.local.remove(key)
    .then(() => true)
    .catch((e: unknown) => {
      // TODO: Some generic error catching here. Could be improved.
      console.error(e);
      return false;
    });
}

const localStorage = {
  getLocalStorage,
  localStorageExists,
  setLocalStorage,
  deleteLocalStorage,
};

export default localStorage;
