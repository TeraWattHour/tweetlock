import { create } from "zustand";

type Vote =
    | {
          votes: number;
          hasVoted: boolean;
          hasError: false;
      }
    | {
          votes: null;
          hasVoted: null;
          hasError: true;
      };

export const useVotesStore = create<{
    votes: {
        [key: string]: Vote;
    };
    setVotesForUser(userId: string, votes: number | null, hasVoted: boolean | null, hasError: boolean): void;
}>((set) => ({
    votes: {},
    setVotesForUser(userId, votes, hasVoted, hasError) {
        set((state) => {
            const snap = state.votes;
            snap[userId] = {
                votes,
                hasVoted,
                hasError
            } as Vote;
            return snap;
        });
    }
}));
