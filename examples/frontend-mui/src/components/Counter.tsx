import Button from "@mui/material/Button";

type Props = {
  count: number;
  increment: () => void;
};
export default function Counter({ count, increment }: Props) {
  return (
    <Button variant="outlined" onClick={increment}>
      Count is {count}
    </Button>
  );
}
