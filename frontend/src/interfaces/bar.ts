export interface BarData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string;
  }[];
}

export interface BarDataOptions {
  responsive: boolean;
  legend: {
    display: boolean;
  };
  title: {
    display: boolean;
    text: string;
    fontSize: number;
    fontColor: string;
  };
  tooltips: {
    backgroundColor: string;
  };
  scales: {
    xAxes: Axis[];
    yAxes: Axis[];
  };
}

interface Axis {
  stacked: boolean;
  ticks?: {
    beginAtZero: boolean;
  };
  gridLines: {
    display: boolean;
  };
}
