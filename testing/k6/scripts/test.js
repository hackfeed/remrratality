import { check } from "k6";
import http from "k6/http";

const baseScenario = {
  executor: "constant-arrival-rate",
  rate: 15000,
  timeUnit: "1s",
  duration: "1m",
  gracefulStop: "0s",
  preAllocatedVUs: 100,
  maxVUs: 300,
};
const token =
  "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbWFpbCI6InRlc3RAdGVzdC5jb20iLCJVc2VySUQiOiI2MTU4YjAzYmQ4Zjc3OTNlMzJjYjQ2ODQiLCJleHAiOjE2MzMyMTQyODZ9.Mr43-hCIKnWS4NmrgWWyqlSDm9aDe11kKcEHyvToYhw";
const filename = "e12a8481-0c2a-4b11-8fad-d574dc41c774.csv";
const period_start = "2019-01-01";
const period_end = "2021-01-01";

export const options = {
  insecureSkipTLSVerify: true,
  systemTags: ["scenario", "check"],
  scenarios: {
    getAnalytics: Object.assign(
      {
        exec: "getAnalytics",
        env: { URL: "https://remrratality.com:8003/api/v1/analytics/mrr" },
      },
      baseScenario
    ),
    getFiles: Object.assign(
      {
        exec: "getFiles",
        env: { URL: "https://remrratality.com:8003/api/v1/files/load" },
        startTime: "1m",
      },
      baseScenario
    ),
  },
};

export const getAnalytics = () => {
  const url = __ENV.URL;
  const params = {
    headers: {
      "Content-Type": "application/json",
      token,
    },
  };
  const createData = JSON.stringify({
    filename,
    period_start,
    period_end,
  });

  const requests = {
    createData: {
      method: "POST",
      url,
      params,
      body: createData,
    },
  };

  const responses = http.batch(requests);
  const createResp = responses.createData;

  check(createResp, {
    "status is 200": (r) => r.status === 200,
  });
};

export const getFiles = () => {
  const url = __ENV.URL;
  const params = {
    headers: {
      token,
    },
  };
  const requests = {
    getData: {
      method: "GET",
      url,
      params,
    },
  };

  const responses = http.batch(requests);
  const getResp = responses.getData;

  check(getResp, {
    "status is 200": (r) => r.status === 200,
  });
};
