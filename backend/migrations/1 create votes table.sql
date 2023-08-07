create table votes (
	target_id varchar(19) not null,
	user_id varchar(36) not null,
		
	unique index ux_vote(target_id, user_id),
	index ix_vote_target(target_id),
	index ix_vote_user(user_id)
);