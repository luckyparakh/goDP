# Skill: revise

Use when the learner says "quiz me", "revise" or "drill". It helps memory, so go
slowly and do not give hints.

Read first: [../rules.md](../rules.md) and `PROGRESS.md`.

## Mode A: quiz

1. Pick 3 questions from the `THEORY.md` files of **finished** patterns (the ones
   in `PROGRESS.md`). Use different patterns, and prefer questions whose
   original `Guess` was wrong.
2. Ask them **one at a time**, with no hints. Wait for each answer.
3. Compare each answer with the learner's earlier `Answer` line and give a
   sharper version if needed.
4. If one pattern keeps going wrong, suggest redoing its delete test.

## Mode B: drill (a real Go file)

1. The learner chooses a Go file. If it may be company or customer code, remind
   them not to put names from it into anything they commit.
2. Ask these three, and **wait for their answers first**:
   1. What varies here, and is that change isolated or spread across the file?
   2. Which pattern is **missing**?
   3. Which pattern is **over-applied**? (An abstraction with one
      implementation and one caller counts.)
3. Then give your own answers and compare. Question 3 matters most for the first
   six months, so spend the most time there.
4. Suggest a log line that uses pattern names only, never file or company names.
   The learner writes it into `PROGRESS.md` themselves.
