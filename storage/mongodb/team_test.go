// Copyright 2017 tsuru authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package mongodb

import (
	"sort"

	"github.com/tsuru/tsuru/storage"
	"gopkg.in/check.v1"
)

var Repo = &TeamRepository{}

func (s *S) TestInsert(c *check.C) {
	t := storage.Team{Name: "teamname", CreatingUser: "me@example.com"}
	err := Repo.Insert(t)
	c.Assert(err, check.IsNil)
	team, err := Repo.FindByName(t.Name)
	c.Assert(err, check.IsNil)
	c.Assert(team.Name, check.Equals, t.Name)
	c.Assert(team.CreatingUser, check.Equals, t.CreatingUser)
}

func (s *S) TestInsertDuplicateTeam(c *check.C) {
	t := storage.Team{Name: "teamname", CreatingUser: "me@example.com"}
	err := Repo.Insert(t)
	c.Assert(err, check.IsNil)
	err = Repo.Insert(t)
	c.Assert(err, check.Equals, storage.ErrTeamAlreadyExists)
}

func (s *S) TestFindAll(c *check.C) {
	err := storage.TeamRepository.Insert(storage.Team{Name: "corrino"})
	c.Assert(err, check.IsNil)
	err = storage.TeamRepository.Insert(storage.Team{Name: "fenring"})
	c.Assert(err, check.IsNil)
	teams, err := Repo.FindAll()
	c.Assert(err, check.IsNil)
	c.Assert(teams, check.HasLen, 2)
	names := []string{teams[0].Name, teams[1].Name}
	sort.Strings(names)
	c.Assert(names, check.DeepEquals, []string{"corrino", "fenring"})
}

func (s *S) TestFindByName(c *check.C) {
	t := storage.Team{Name: "myteam"}
	err := Repo.Insert(t)
	c.Assert(err, check.IsNil)
	team, err := Repo.FindByName(t.Name)
	c.Assert(err, check.IsNil)
	c.Assert(team.Name, check.Equals, t.Name)
}

func (s *S) TestFindByNameNotFound(c *check.C) {
	team, err := Repo.FindByName("wat")
	c.Assert(err, check.Equals, storage.ErrTeamNotFound)
	c.Assert(team, check.IsNil)
}

func (s *S) TestFindByNames(c *check.C) {
	t1 := storage.Team{Name: "team1"}
	err := Repo.Insert(t1)
	c.Assert(err, check.IsNil)
	t2 := storage.Team{Name: "team2"}
	err = Repo.Insert(t2)
	c.Assert(err, check.IsNil)
	t3 := storage.Team{Name: "team3"}
	err = Repo.Insert(t3)
	c.Assert(err, check.IsNil)
	teams, err := Repo.FindByNames([]string{t1.Name, t2.Name, "unknown"})
	c.Assert(err, check.IsNil)
	c.Assert(teams, check.DeepEquals, []storage.Team{t1, t2})
}

func (s *S) TestFindByNamesNotFound(c *check.C) {
	t1 := storage.Team{Name: "team1"}
	err := Repo.Insert(t1)
	c.Assert(err, check.IsNil)
	teams, err := Repo.FindByNames([]string{"unknown", "otherteam"})
	c.Assert(err, check.IsNil)
	c.Assert(teams, check.HasLen, 0)
}

func (s *S) TestDelete(c *check.C) {
	team := storage.Team{Name: "atreides"}
	err := Repo.Insert(team)
	c.Assert(err, check.IsNil)
	err = Repo.Delete(team)
	c.Assert(err, check.IsNil)
	t, err := Repo.FindByName("atreides")
	c.Assert(err, check.Equals, storage.ErrTeamNotFound)
	c.Assert(t, check.IsNil)
}

func (s *S) TestDeleteTeamNotFound(c *check.C) {
	err := Repo.Delete(storage.Team{Name: "myteam"})
	c.Assert(err, check.Equals, storage.ErrTeamNotFound)
}
