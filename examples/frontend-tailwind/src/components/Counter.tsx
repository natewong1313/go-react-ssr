type Props = {
  count: number;
  increment: () => void;
};
export default function Counter({ count, increment }: Props) {
  return (
    <button onClick={increment} className="text-white bg-zinc-900 py-2.5 px-4 border-2 border-transparent hover:border-sky-500 rounded-md">
      Count is {count}
    </button>
  );
}
