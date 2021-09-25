export interface BarData {
  labels: string[];
  datasets: {
    label: string;
    data: number[];
    backgroundColor: string;
  }[];
}

export interface BarOptions {
  responsive: boolean;
  plugins: {
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
  };
  scales: {
    x: Axis;
    y: Axis;
  };
}

interface Axis {
  stacked: boolean;
  ticks?: {
    beginAtZero: boolean;
  };
  grid: {
    display: boolean;
  };
}
