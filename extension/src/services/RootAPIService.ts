import RootAPIResponse from '../models/RootAPIResponse.model';
import { getBaseURL } from '../utils/urls';

async function getName(): Promise<string> {
  const response = await fetch(getBaseURL(), {
    method: 'GET',
  });
  if (response.status !== 200) {
    return '';
  }
  const jsonData = await response.json() as RootAPIResponse;
  return jsonData.name;
}

const RootAPIService = {
  getName,
};

export default RootAPIService;
