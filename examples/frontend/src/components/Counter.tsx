import styles from "./Counter.module.css";

type Props = {
  count: number;
  increment: () => void;
};
export default function Counter({ count, increment }: Props) {
  return (
    <button className={styles.counter} onClick={increment}>
      Count is {count}
    </button>
  );
}
